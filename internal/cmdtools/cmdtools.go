package cmdtools

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/edwarnicke/exechelper"
)

func RunCmd(command string, override ...*exechelper.Option) error {
	var options []*exechelper.Option

	if len(override) > 0 {
		options = override
	} else {
		options = []*exechelper.Option{exechelper.WithStdout(os.Stdout), exechelper.WithStderr(os.Stderr)}
	}

	output, err := RunCmdAndGetOutput(command, options...)

	if err != nil {
		return err
	} else if output.ErrorCode != 0 {
		return fmt.Errorf("Failed to run command [%s], error code: %d, stderr: %s",
			command, output.ErrorCode, string(output.Stderr))
	}

	return nil
}

type CmdOutput struct {
	Stdout    []byte
	Stderr    []byte
	ErrorCode int
}

// The behavior of the func is similar to the one of golang http client.
// Getting error code from the command is not considered an error, only
// a generic error when we failed to get any error code.
// It's up to handling code to interpret the error code value. If you want
// to treat all non zero error code as errors, make a wrapper function
func RunCmdAndGetOutput(command string, additional ...*exechelper.Option) (*CmdOutput, error) {
	var outputBuffer bytes.Buffer
	var errBuffer bytes.Buffer

	options := []*exechelper.Option{
		exechelper.WithStdout(&outputBuffer),
		exechelper.WithStderr(&errBuffer),
	}
	options = append(options, additional...)

	err := exechelper.Run(command, options...)

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return &CmdOutput{
				Stdout:    outputBuffer.Bytes(),
				Stderr:    errBuffer.Bytes(),
				ErrorCode: exitError.ExitCode(),
			}, nil
		}

		return nil, err
	}

	return &CmdOutput{
		Stdout: outputBuffer.Bytes(),
		Stderr: errBuffer.Bytes(),
	}, nil

}

func RunCmdAndGetOutputNoErrCode(command string, additional ...*exechelper.Option) (*CmdOutput, error) {
	output, err := RunCmdAndGetOutput(command, additional...)

	if err != nil {
		return nil, err
	} else if output.ErrorCode != 0 {
		return nil, fmt.Errorf("Failed to run command [%s], error code: %d, stderr: %s",
			command, output.ErrorCode, string(output.Stderr))
	}

	return output, nil
}
