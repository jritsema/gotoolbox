package gotoolbox

import (
	"bytes"
	"os/exec"
)

//ExecCmd executes a command and waits
func ExecCmd(cmd *exec.Cmd) error {

	var outbuf, errbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Start()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}
