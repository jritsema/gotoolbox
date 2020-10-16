package gotoolbox

import (
	"bytes"
	"os"
	"os/exec"
)

//GetEnvWithDefault returns the value of an enviornment variable
//or a default value if the environment variables is not set
func GetEnvWithDefault(key, defaultValue string) string {
	result := os.Getenv(key)
	if result == "" {
		return defaultValue
	}
	return result
}

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
