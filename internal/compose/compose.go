// Package compose wraps `docker compose` invocations against a specific
// project directory, compose file, and env file — the shape every command
// in the springboot wizard needs to run against a docker-compose.prod.yml
// project like Kuniel's.
package compose

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Project identifies a docker-compose project on disk.
type Project struct {
	Dir     string // directory containing the compose file (cwd for the command)
	File    string // e.g. "docker-compose.prod.yml"
	EnvFile string // e.g. ".env.prod"
}

// Args builds the full "docker compose ..." argument list for extra —
// pulled out from command() so the construction logic is testable without
// actually invoking docker.
func (p Project) Args(extra ...string) []string {
	base := []string{"compose", "-f", p.File, "--env-file", p.EnvFile}
	return append(base, extra...)
}

func (p Project) command(args ...string) *exec.Cmd {
	cmd := exec.Command("docker", p.Args(args...)...)
	cmd.Dir = p.Dir
	return cmd
}

// Run streams a compose command's stdout/stderr to the current process —
// use for anything the operator should watch live (up, pull, logs -f).
func (p Project) Run(args ...string) error {
	cmd := p.command(args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Output runs a compose command and returns its combined, trimmed output.
func (p Project) Output(args ...string) (string, error) {
	out, err := p.command(args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// ContainerID returns the container ID for a compose service, including one
// that has already exited or is crash-looping. Stopped containers count on
// purpose: the reason to look up an ID at all is usually to recover what a
// *failed* container was started from, and -q alone would skip exactly that
// case. Compose keeps one container per service unless it was scaled, so the
// first line is normally the only line.
func (p Project) ContainerID(service string) (string, error) {
	out, err := p.Output("ps", "-aq", service)
	if err != nil {
		return "", fmt.Errorf("compose ps -aq %s: %w\n%s", service, err, out)
	}
	id := strings.SplitN(out, "\n", 2)[0]
	if id == "" {
		return "", fmt.Errorf("service %q has no container", service)
	}
	return id, nil
}

// ImageID returns the image ID (sha256:...) a container was started from —
// stable even after a new image is pulled and a floating tag (e.g. :latest)
// is repointed, which is what makes rollback possible.
func ImageID(containerID string) (string, error) {
	out, err := exec.Command("docker", "inspect", "--format", "{{.Image}}", containerID).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker inspect %s: %w\n%s", containerID, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// LocalImageID returns the image ID that a local image reference currently
// resolves to. Comparing it against the ID a container was started from is
// what separates "the pull brought a new build" from "the tag hasn't moved" —
// and in the latter case there is no earlier image to roll back to.
func LocalImageID(imageRef string) (string, error) {
	out, err := exec.Command("docker", "image", "inspect", "--format", "{{.Id}}", imageRef).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker image inspect %s: %w\n%s", imageRef, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// UserSpec returns an image's configured USER verbatim: "" for an image that
// runs as root, or whatever the Dockerfile set ("1001", "appuser",
// "appuser:appgroup").
func UserSpec(imageRef string) (string, error) {
	out, err := exec.Command("docker", "image", "inspect", "--format", "{{.Config.User}}", imageRef).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("docker image inspect %s: %w\n%s", imageRef, err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

// ImageUserIDs reports the numeric uid/gid a container started from imageRef
// runs as — the numbers a host-side chown needs to hand a bind-mounted file to
// that container. USER is usually a name, and only the image's own passwd file
// knows what it resolves to, so the numbers are read by running id(1) inside
// it rather than guessed from the inspect output.
func ImageUserIDs(imageRef string) (int, int, error) {
	spec, err := UserSpec(imageRef)
	if err != nil {
		return 0, 0, err
	}
	if isRootSpec(spec) {
		return 0, 0, nil
	}

	uid, uidErr := idInImage(imageRef, "-u")
	gid, gidErr := idInImage(imageRef, "-g")
	if uidErr == nil && gidErr == nil {
		return uid, gid, nil
	}
	// Distroless-style images ship no id(1). A numeric USER still tells us
	// everything we need; a named one there is genuinely unresolvable.
	if specUID, specGID, ok := parseNumericUserSpec(spec); ok {
		return specUID, specGID, nil
	}
	if uidErr != nil {
		return 0, 0, fmt.Errorf("resolve USER %q of %s: %w", spec, imageRef, uidErr)
	}
	return 0, 0, fmt.Errorf("resolve USER %q of %s: %w", spec, imageRef, gidErr)
}

// idInImage runs id(1) as the image's own USER. -u and -g are asked for one at
// a time because bare id(1) answers with "uid=1001(appuser) gid=1001(appgroup)
// groups=..." — a single number per call needs no parsing to go wrong.
func idInImage(imageRef, flag string) (int, error) {
	cmd := exec.Command("docker", "run", "--rm", "--network", "none",
		"--entrypoint", "id", imageRef, flag)
	// Stdout only: docker writes warnings (a platform mismatch, say) to stderr,
	// and a combined stream would leave the number we came for unparseable.
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("docker run --entrypoint id %s: %w\n%s", imageRef, err, stderr.String())
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, fmt.Errorf("parse `id %s` output %q: %w", flag, out, err)
	}
	return n, nil
}

// isRootSpec covers every way an image can say "root", including the empty
// spec Docker reports when the Dockerfile never set USER at all.
func isRootSpec(spec string) bool {
	switch spec {
	case "", "0", "root", "0:0", "root:root":
		return true
	}
	return false
}

// parseNumericUserSpec pulls uid/gid out of a numeric spec ("1001",
// "1001:1001"). A spec naming a user or group instead can't be resolved
// without the image's passwd file, and reports ok=false.
func parseNumericUserSpec(spec string) (int, int, bool) {
	user, group, hasGroup := strings.Cut(spec, ":")
	uid, err := strconv.Atoi(user)
	if err != nil {
		return 0, 0, false
	}
	if !hasGroup {
		// The container's real gid comes from the image's passwd entry for
		// this uid, which is exactly what we couldn't read. Mirroring the uid
		// costs nothing: the files stay 0700/0600, where only the owner bits
		// decide anything.
		return uid, uid, true
	}
	gid, err := strconv.Atoi(group)
	if err != nil {
		return uid, uid, true
	}
	return uid, gid, true
}

// Tag re-points imageRef (e.g. "ghcr.io/org/app:latest") at imageID — used to
// roll back by making the next `up` pick up the old image without re-pulling.
func Tag(imageID, imageRef string) error {
	out, err := exec.Command("docker", "tag", imageID, imageRef).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker tag %s %s: %w\n%s", imageID, imageRef, err, out)
	}
	return nil
}
