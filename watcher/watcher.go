package watcher

import (
	"os"
	"time"

	"github.com/srivatsa-bot/gomon/logger"
)

type FileWatcher struct {
	filename    string
	lastModTime time.Time
	serverproc  *ServerProcess
}

// to get intial stat
func NewFileWatcher(filename string) (*FileWatcher, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}
	//return the filewatcher struct
	return &FileWatcher{
		filename:    filename,
		lastModTime: fileInfo.ModTime(),
	}, nil

}

// to start the server using the type filewatcher's serverproc
// method for file watcher struct and start the server
func (fw *FileWatcher) Start() error {
	var err error //pointer is used so no :=(declares again) serveproc is already declared in type definintion
	fw.serverproc, err = StartServer(fw.filename)
	if err != nil {
		return err
	}

	return nil
}

// watcher logic
func (fw *FileWatcher) Watch() error {
	for {
		time.Sleep(1 * time.Second)

		fileinfo, err := os.Stat(fw.filename)
		if err != nil {
			logger.Error("File doesn't exist: %v", err)
			continue
		}
		if fileinfo.ModTime().After(fw.lastModTime) {
			fw.lastModTime = fileinfo.ModTime()
			logger.Info("Detected changes in %s - restarting the server", fw.filename)

			//kill the process
			if err := fw.serverproc.Kill(); err != nil {
				logger.Error("Error killing the process: %v", err)
			}

			//restart the process again
			if err := fw.Start(); err != nil {
				logger.Error("Failed to restart the server: %v", err)
				return err
			}
		}

	}
}

// cleanup code when os signal is recieved
func (fw *FileWatcher) Cleanup() {
	//kill if a process is running
	if fw.serverproc != nil {
		logger.Info("Cleaning up server process...")
		fw.serverproc.Kill()
		logger.Info("Server process terminated")
	}
}
