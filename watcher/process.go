package watcher

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/srivatsa-bot/gomon/logger"
)

type ServerProcess struct {
	cmd *exec.Cmd
}

// function to start server
func StartServer(filename string) (*ServerProcess, error) {
	//handling multiple files
	var cmd *exec.Cmd
	// Get file extension
	ext := strings.ToLower(filepath.Ext(filename))

	//support for interpreted lang(compiled will be updated in future)
	switch ext {
	case ".go":
		cmd = exec.Command("go", "run", filename)
	case ".js":
		cmd = exec.Command("node", filename)
	case ".py":
		cmd = exec.Command("python3", filename)
	default:
		// For unknown file types, try to execute directly if executable
		logger.Info("File type not supported yet")
		os.Exit(1)
	}

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
		Setpgid: true,
		// Pdeathsig: syscall.SIGTERM, not needed as we are handling various signlas in kill() method
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
	goos := runtime.GOOS //get the os of the system

	//get pgid(for unix)
	pgid, err := syscall.Getpgid(sp.cmd.Process.Pid)
	//kill process gracefully using sigterm
	if err == nil {

		switch goos {
		case "linux":
			if err = syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
				syscall.Kill(-pgid, syscall.SIGKILL) //if sigterm fails force kill it using sigkill
			}
		case "darwin":
			if err = syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
				syscall.Kill(-pgid, syscall.SIGKILL) 
			}

		}

	}

	//wait for process to exit gracefully
	sp.cmd.Wait()
	time.Sleep(100 * time.Millisecond)
	return nil
}
