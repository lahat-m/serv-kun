package system

import (
	"os"
	"os/exec"
)

// RunCmd runs name with args, streaming stdout/stderr to the current process.
func RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
