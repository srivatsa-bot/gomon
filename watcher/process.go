package watcher

import (
	"os/exec"
	"syscall"
	"time"

	"github.com/srivatsa-bot/gomon/logger"
)

type ServerProcess struct {
	cmd *exec.Cmd
}

// function to start server
func StartServer(filename string) (*ServerProcess, error) {
	cmd := exec.Command("go", "run", filename)

	//pipe output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	//allocate process to process group
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGTERM,
	}

	//start the command
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	//code to log starts
	go logger.LogOutput(stdout, "Server")
	go logger.LogOutput(stderr, "Server Error")
	logger.Info("Server started with PID: %d", cmd.Process.Pid)
	//code to log ends

	return &ServerProcess{cmd: cmd}, nil
}

//function to kill process

func (sp *ServerProcess) Kill() error {
	if sp.cmd == nil || sp.cmd.Process == nil {
		return nil
	}

	//get pgid
	pgid, err := syscall.Getpgid(sp.cmd.Process.Pid)
	//kill process gracefully using sigterm
	if err == nil {
		err = syscall.Kill(-pgid, syscall.SIGTERM)
		if err != nil {
			syscall.Kill(-pgid, syscall.SIGKILL) //if sigterm fails force kill it using sigkill
		}
	}

	//wait for process to exit gracefully
	sp.cmd.Wait()
	time.Sleep(100 * time.Millisecond)
	return nil
}
