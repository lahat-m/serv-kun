// Package compose wraps `docker compose` invocations against a specific
// project directory, compose file, and env file — the shape every command
// in the springboot wizard needs to run against a docker-compose.prod.yml
// project like Kuniel's.
package compose

import (
	"fmt"
	"os"
	"os/exec"
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

// Tag re-points imageRef (e.g. "ghcr.io/org/app:latest") at imageID — used to
// roll back by making the next `up` pick up the old image without re-pulling.
func Tag(imageID, imageRef string) error {
	out, err := exec.Command("docker", "tag", imageID, imageRef).CombinedOutput()
	if err != nil {
		return fmt.Errorf("docker tag %s %s: %w\n%s", imageID, imageRef, err, out)
	}
	return nil
}
