package tests

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func getCmd(t *testing.T, command string) *exec.Cmd {

	//split command into command and args
	split := strings.Split(command, " ")
	require.True(t, len(split) > 0, "no command provided")

	var cmd *exec.Cmd
	if len(split) == 1 {
		cmd = exec.Command(split[0])
	} else {
		cmd = exec.Command(split[0], split[1:]...)
	}
	return cmd
}

// Execute the command, return standard output and error, try a few times if requested
func ExecuteT(t *testing.T, command string) (out string) {
	cmd := getCmd(t, command)
	//bz, err := cmd.CombinedOutput()
	//require.NoError(t, err, string(bz))

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			txt := scanner.Text()
			out += scanner.Text() + "\n"
			fmt.Printf("debug out | %s\n", txt)
		}
	}()
	err = cmd.Start()
	require.NoError(t, err, "Error starting Cmd")
	err = cmd.Wait()
	require.NoError(t, err, "Error waiting for Cmd")

	out = strings.Trim(out, "\n") //trim any new lines
	return out
}

// Asynchronously execute the command, return standard output and error
func GoExecuteT(t *testing.T, command string) (cmd *exec.Cmd, pipeIn io.WriteCloser, pipeOut io.ReadCloser) {
	cmd = getCmd(t, command)
	pipeIn, err := cmd.StdinPipe()
	require.NoError(t, err)
	pipeOut, err = cmd.StdoutPipe()
	require.NoError(t, err)
	cmd.Start()
	return cmd, pipeIn, pipeOut
}
