//go:build docker

// These tests need a working docker daemon and are excluded from the ordinary
// `go test ./...`. Run them with:
//
//	go test -tags docker ./internal/compose/
//
// They exist because the interesting failures in this file aren't logic — they
// are assumptions about how docker reports and applies an image's USER, which
// only a real daemon can confirm.
package compose

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/lahat-m/serv-kun/internal/jwtkeys"
)

// namedUserImage mirrors the shape of a Spring Boot production image: a system
// account created at a fixed uid but referred to by name, which is what makes
// Config.User useless to chown on its own.
const namedUserImage = `FROM debian:bookworm-slim
RUN groupadd --system --gid 1001 appgroup \
 && useradd --system --uid 1001 --gid appgroup --no-create-home appuser
USER appuser
CMD ["java", "-jar", "application.jar"]
`

// renumberedUserImage is the same application image after a base-image bump
// moved its system account to a different uid — the case that breaks a deploy
// whose certs/ was chowned once, at setup, and never revisited.
const renumberedUserImage = `FROM debian:bookworm-slim
RUN groupadd --system --gid 1002 appgroup \
 && useradd --system --uid 1002 --gid appgroup --no-create-home appuser
USER appuser
CMD ["java", "-jar", "application.jar"]
`

// numericNoIDImage is the distroless case: nothing to exec, so the uid has to
// come out of the USER spec itself.
const numericNoIDImage = `FROM scratch
USER 1001:1002
`

// namedNoIDImage is the one genuinely unresolvable combination — a name, and
// no id(1) to resolve it with.
const namedNoIDImage = `FROM scratch
USER appuser
`

func buildImage(t *testing.T, dockerfile, tag string) string {
	t.Helper()

	cmd := exec.Command("docker", "build", "-q", "-t", tag, "-")
	cmd.Stdin = strings.NewReader(dockerfile)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("docker build %s: %v\n%s", tag, err, out)
	}
	t.Cleanup(func() { _ = exec.Command("docker", "image", "rm", "-f", tag).Run() })
	return tag
}

func TestImageUserIDsResolvesNamedUser(t *testing.T) {
	image := buildImage(t, namedUserImage, "servkun-test-named")

	// The premise of the whole exercise: what docker reports is a name.
	spec, err := UserSpec(image)
	if err != nil {
		t.Fatalf("UserSpec: %v", err)
	}
	if spec != "appuser" {
		t.Fatalf("UserSpec = %q, want %q", spec, "appuser")
	}

	uid, gid, err := ImageUserIDs(image)
	if err != nil {
		t.Fatalf("ImageUserIDs: %v", err)
	}
	if uid != 1001 || gid != 1001 {
		t.Errorf("ImageUserIDs = %d:%d, want 1001:1001", uid, gid)
	}
}

// An image that never set USER runs as root and needs no chown at all — and
// must not cost a container start to find that out.
func TestImageUserIDsRootImage(t *testing.T) {
	uid, gid, err := ImageUserIDs("debian:bookworm-slim")
	if err != nil {
		t.Fatalf("ImageUserIDs: %v", err)
	}
	if uid != 0 || gid != 0 {
		t.Errorf("ImageUserIDs = %d:%d, want 0:0", uid, gid)
	}
}

// No id(1) to run, but the numbers are right there in the spec.
func TestImageUserIDsFallsBackToNumericSpec(t *testing.T) {
	image := buildImage(t, numericNoIDImage, "servkun-test-noid-numeric")

	uid, gid, err := ImageUserIDs(image)
	if err != nil {
		t.Fatalf("ImageUserIDs: %v", err)
	}
	if uid != 1001 || gid != 1002 {
		t.Errorf("ImageUserIDs = %d:%d, want 1001:1002", uid, gid)
	}
}

// Neither route works here, and reporting 0:0 would be worse than failing:
// setup would silently leave the keys root-only and the app would crash-loop
// with no explanation of why.
func TestImageUserIDsFailsOnUnresolvableName(t *testing.T) {
	image := buildImage(t, namedNoIDImage, "servkun-test-noid-named")

	uid, gid, err := ImageUserIDs(image)
	if err == nil {
		t.Fatalf("ImageUserIDs = %d:%d, nil — want an error", uid, gid)
	}
	if !strings.Contains(err.Error(), "appuser") {
		t.Errorf("error %q doesn't name the USER it failed on", err)
	}
}

// The end-to-end case, and the bug this all exists for: keys written root-only
// into a volume the app mounts, read back by a container running as the image's
// non-root USER. Before GrantUser that read fails on the 0700 directory even
// though public.pem itself is world-readable — which is exactly how it showed
// up in production, as a permission error on the *public* key.
func TestGrantUserMakesKeysReadableInContainer(t *testing.T) {
	certsDir, volume := certsFixture(t)
	image := buildImage(t, namedUserImage, "servkun-test-named-e2e")

	readKey := func() ([]byte, error) { return readKeyAs(image, volume) }

	out, err := readKey()
	if err == nil {
		t.Fatalf("the app user could already read the key before GrantUser — the fixture doesn't reproduce the bug\n%s", out)
	}
	if !strings.Contains(string(out), "Permission denied") {
		t.Fatalf("read failed for the wrong reason: %s", out)
	}

	uid, gid, err := ImageUserIDs(image)
	if err != nil {
		t.Fatalf("ImageUserIDs: %v", err)
	}
	if err := jwtkeys.GrantUser(certsDir, uid, gid); err != nil {
		t.Fatalf("GrantUser: %v", err)
	}

	if out, err := readKey(); err != nil {
		t.Fatalf("the app user still cannot read the key after GrantUser: %v\n%s", err, out)
	}

	// The point of chowning rather than chmodding: one container user gained
	// access, the rest of the host didn't.
	info, err := os.Stat(certsDir + "/private.pem")
	if err != nil {
		t.Fatalf("Stat: %v", err)
	}
	if got := info.Mode().Perm(); got != 0600 {
		t.Errorf("private.pem mode = %o after GrantUser, want 600", got)
	}
}

// The deploy case: an image whose user was renumbered can't read keys chowned
// for the image before it, and — the part that matters — rolling back to that
// earlier image can't read keys chowned for the new one either. Ownership has
// to be reasserted in both directions, or the rollback that exists to survive
// a bad deploy fails alongside it.
func TestGrantUserFollowsAUIDChangeInBothDirections(t *testing.T) {
	certsDir, volume := certsFixture(t)
	oldImage := buildImage(t, namedUserImage, "servkun-test-uid-1001")
	newImage := buildImage(t, renumberedUserImage, "servkun-test-uid-1002")

	grantTo := func(image string) {
		t.Helper()
		uid, gid, err := ImageUserIDs(image)
		if err != nil {
			t.Fatalf("ImageUserIDs(%s): %v", image, err)
		}
		if err := jwtkeys.GrantUser(certsDir, uid, gid); err != nil {
			t.Fatalf("GrantUser(%d:%d): %v", uid, gid, err)
		}
	}

	// Deploy: keys are handed to the image that's arriving.
	grantTo(oldImage)
	if out, err := readKeyAs(oldImage, volume); err != nil {
		t.Fatalf("uid 1001 cannot read keys granted to it: %v\n%s", err, out)
	}
	if _, err := readKeyAs(newImage, volume); err == nil {
		t.Fatal("uid 1002 read keys owned by 1001 — the fixture isn't reproducing a uid change")
	}

	grantTo(newImage)
	if out, err := readKeyAs(newImage, volume); err != nil {
		t.Fatalf("uid 1002 cannot read keys after the deploy re-grant: %v\n%s", err, out)
	}

	// Rollback: the old image is restored, and the keys have to follow it back.
	if _, err := readKeyAs(oldImage, volume); err == nil {
		t.Fatal("uid 1001 read keys owned by 1002 — rollback wouldn't need a re-grant")
	}
	grantTo(oldImage)
	if out, err := readKeyAs(oldImage, volume); err != nil {
		t.Fatalf("rollback re-grant left uid 1001 unable to read its keys: %v\n%s", err, out)
	}
}

// certsFixture lays down a keypair exactly as jwtkeys.Generate leaves one:
// root-owned, directory 0700, private key 0600. It needs a docker volume that
// is both mounted here and mountable by the daemon, so the container side of
// the test sees the same files with the same ownership.
func certsFixture(t *testing.T) (certsDir, volume string) {
	t.Helper()

	certsDir = os.Getenv("SERVKUN_TEST_CERTS_DIR")
	volume = os.Getenv("SERVKUN_TEST_CERTS_VOLUME")
	if certsDir == "" || volume == "" {
		t.Skip("set SERVKUN_TEST_CERTS_DIR and SERVKUN_TEST_CERTS_VOLUME to the same docker volume, mounted here and available to the daemon")
	}

	if err := os.MkdirAll(certsDir, 0700); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.Chmod(certsDir, 0700); err != nil {
		t.Fatalf("Chmod: %v", err)
	}
	for name, mode := range map[string]os.FileMode{"private.pem": 0600, "public.pem": 0644} {
		if err := os.WriteFile(certsDir+"/"+name, []byte("k\n"), mode); err != nil {
			t.Fatalf("WriteFile %s: %v", name, err)
		}
		if err := os.Chmod(certsDir+"/"+name, mode); err != nil {
			t.Fatalf("Chmod %s: %v", name, err)
		}
	}
	if err := os.Chown(certsDir, 0, 0); err != nil {
		t.Fatalf("Chown to root: %v", err)
	}
	return certsDir, volume
}

// readKeyAs reads the private key from inside a container running as image's
// own USER — the only view of these permissions that decides anything.
func readKeyAs(image, volume string) ([]byte, error) {
	return exec.Command("docker", "run", "--rm", "--network", "none",
		"-v", volume+":/certs", "--entrypoint", "cat", image, "/certs/private.pem").CombinedOutput()
}
