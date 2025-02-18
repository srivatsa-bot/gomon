package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/srivatsa-bot/gomon/logger"
	"github.com/srivatsa-bot/gomon/watcher"
)

func main() {
	//siganl handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,  //ctrl+c
		syscall.SIGTERM, //termination signal from os
		syscall.SIGQUIT, //ctrl+\
		syscall.SIGHUP,  // closing terminal
	)

	//get the file stat
	fw, err := watcher.NewFileWatcher("server.go")
	if err != nil {
		logger.Error("Failed to create file watcher: %v", err)
		return
	}
	//start the server
	if err = fw.Start(); err != nil {
		logger.Error("Failed to start server: %v", err)
	}

	//start the sever first cuz, serverproc is atteched after the start method
	//then do defer and handle signal cuz both use fw.cleanup method which takes fw.serverproc as inputs
	//this gurantess fw is properly initialized

	//clean the process
	defer fw.Cleanup()

	//handling interrupt in a seperate go routine
	go func() {
		sig := <-sigChan
		logger.Info("\nReceived %v signal. Shutting down...", sig)
		fw.Cleanup()
		os.Exit(0)

	}()

	//watch the file for changes
	if err := fw.Watch(); err != nil {
		logger.Error("Watcher error: %v", err)
	}

}
