package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/srivatsa-bot/gomon/logger"
	"github.com/srivatsa-bot/gomon/watcher"
)

var Version string = "1.0.0"

func main() {

	//cli kit
	fileName := flag.String("file", "server.go", "File to watch and run(Handles only one file)")
	showVersion := flag.Bool("version", false, "Show version information")
	showHelp := flag.Bool("help", false, "Show help information")

	flag.Parse() //takes input from cli and assigns it flags

	// Show version and exit
	if *showVersion {
		fmt.Printf("gomon version %s\n", Version)
		return
	}
	// Show help and exit
	if *showHelp {
		fmt.Println("Usage of gomon:")
		fmt.Println("  gomon [flags] [file]")
		fmt.Println("\nFlags:")
		flag.PrintDefaults() //prints all the default flags
		return
	}

	//get filename form tag
	targetFile := *fileName
	//to make gomon server.go(nonflag argument work)
	if flag.NArg() > 0 {
		targetFile = flag.Arg(0) //first non flag argument

		//if user enters more than one non flag argument
		if flag.NArg() > 1 {
			logger.Info("Warning: Multiple files provided. Only watching %s", *fileName)
		}
	}

	//check if file exists then proceed
	if _, err := os.Stat(targetFile); os.IsNotExist(err) { //os.isnotexist returns true if file doesnt exist
		fmt.Printf("File %s does not exist\n", targetFile)
		os.Exit(1)
	}

	//siganl handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGINT,  //ctrl+c
		syscall.SIGTERM, //termination signal from os
		syscall.SIGQUIT, //ctrl+\
		syscall.SIGHUP,  // closing terminal
	)

	//get the file stat
	fw, err := watcher.NewFileWatcher(targetFile)
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
		os.Exit(0) //exit sucesfully

	}()

	//watch the file for changes
	if err := fw.Watch(); err != nil {
		logger.Error("Watcher error: %v", err)
	}

}
