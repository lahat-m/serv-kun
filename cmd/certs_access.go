package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/lahat-m/serv-kun/internal/compose"
	"github.com/lahat-m/serv-kun/internal/jwtkeys"
)

// appImageRef returns the image to interrogate for the app service. The
// container knows exactly which image it was created from, which is the honest
// answer when the compose file builds locally instead of pulling the tag, or
// when a floating tag has since been repointed. The tag is the fallback for a
// container that doesn't exist yet.
func appImageRef(proj compose.Project, tag string) string {
	containerID, err := proj.ContainerID("app")
	if err != nil {
		return tag
	}
	imageID, err := compose.ImageID(containerID)
	if err != nil {
		return tag
	}
	return imageID
}

// ensureCertsReadable hands certs/ to the account imageRef's USER resolves to,
// so a non-root container can read the JWT keypair mounted into it. Keys are
// generated root-owned, which such an image can see but not open — it then
// dies at startup on its own signing key.
//
// This is asserted against whichever image is about to run rather than once at
// setup, because the answer changes: an image that renumbers its user, or a
// rollback to an image that used a different one, leaves certs/ owned by a uid
// that no longer exists in the container. Ownership is cheap to reassert and
// the cost of getting it wrong is a crash-looping app.
//
// Reports whether ownership was actually handed over. Nothing here is fatal —
// a failure to resolve the user is worth a printed chown, not an aborted
// deploy that leaves the stack half-moved.
func ensureCertsReadable(proj compose.Project, imageRef string) bool {
	certsDir := filepath.Join(proj.Dir, "certs")
	if !jwtkeys.Exists(certsDir) {
		// No mounted keypair, nothing to hand over — and no reason to spend
		// two container starts finding that out.
		return false
	}

	uid, gid, err := compose.ImageUserIDs(imageRef)
	if err != nil {
		fmt.Printf("! could not read the app image's USER: %v\n", err)
		fmt.Printf("  If it runs as a non-root user it cannot read certs/, and will fail to\n")
		fmt.Printf("  start. Fix with: sudo chown -R <uid>:<gid> %s\n", certsDir)
		return false
	}
	// root reads the keys whatever they're owned by, including keys still
	// owned by the non-root user of a previous image.
	if uid == 0 {
		return false
	}

	if err := jwtkeys.GrantUser(certsDir, uid, gid); err != nil {
		fmt.Printf("! could not hand certs/ to the app user: %v\n", err)
		fmt.Printf("  Fix with: sudo chown -R %d:%d %s\n", uid, gid, certsDir)
		return false
	}
	fmt.Printf("✓ certs/ handed to %d:%d — the app image's USER\n", uid, gid)
	return true
}
