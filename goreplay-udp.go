package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var closeCh chan int

func main() {

	// add line number to log
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if !flag.Parsed() {
		flag.Parse()
	}

	InitPlugins()

	if len(Plugins.Inputs) == 0 || len(Plugins.Outputs) == 0 {
		log.Fatal("Required at least 1 input and 1 output")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		finalize()
		os.Exit(1)
	}()

	if Settings.exitAfter > 0 {
		log.Println("Running gor for a duration of", Settings.exitAfter)
		closeCh = make(chan int)

		time.AfterFunc(Settings.exitAfter, func() {
			log.Println("Stopping gor after", Settings.exitAfter)
			close(closeCh)
		})
	}

	Start(closeCh)

}

func finalize() {
	for _, p := range Plugins.All {
		if cp, ok := p.(io.Closer); ok {
			cp.Close()
		}
	}
}
