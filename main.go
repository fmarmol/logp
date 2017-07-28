package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "impossible to stat stdin: %v\n", err)
		os.Exit(1)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Fprintf(os.Stderr, "Stdin is not in pipe mode, ther is nothng to log\n")
		os.Exit(1)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	strings := make(chan string)
	go func() {
		defer close(strings)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			strings <- s.Text()
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan input: %v\n", err)
			os.Exit(1)
		}
	}()

	finish := make(chan struct{})
	go func() {
		for str := range strings {
			log.Print(str)
		}
		finish <- struct{}{}
	}()

loop:
	for {
		select {
		case <-finish:
			break loop
		case sig := <-signals:
			fmt.Fprintf(os.Stderr, "Program exits after recieved %v signal\n", sig)
			break loop
		}
	}
}
