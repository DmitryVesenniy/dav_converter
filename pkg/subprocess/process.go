package subprocess

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type SysProcess struct {
	OS      string
	Command Command
}

type Command struct {
	NameCommand string
	Args        []string
}

func (sp *SysProcess) Run() error {
	cmd := exec.Command(sp.Command.NameCommand, sp.Command.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error run system process: %w", err)
	}
	cmd.Wait()
	return nil
}

func NewSysProcess(commands map[string]Command) (*SysProcess, error) {
	osSystem := runtime.GOOS

	_command, ok := commands[osSystem]
	if !ok {
		return &SysProcess{}, fmt.Errorf("not found os command")
	}

	return &SysProcess{
		OS:      osSystem,
		Command: _command,
	}, nil
}
