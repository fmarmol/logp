package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"time"
)

var layout = "2006/01/02 15:03:04"

// Msg is a struct to wrap stdin input
type Msg struct {
	time time.Time
	text string
}

func log(msg Msg, previousTime *time.Time) {
	if (*previousTime).Unix() != (time.Time{}).Unix() {
		fmt.Printf("%v %v. Time elapsed since last message: %v\n", msg.time.Format(layout), msg.text, msg.time.Sub(*previousTime))
	} else {
		fmt.Printf("%v %v\n", msg.time.Format(layout), msg.text)
	}
	*previousTime = msg.time
}

func main() {
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "impossible to stat stdin: %v\n", err)
		os.Exit(1)
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		fmt.Fprintf(os.Stderr, "Stdin is not in pipe mode, ther is nothing to log\n")
		os.Exit(1)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)

	msgs := make(chan Msg)
	go func() {
		defer close(msgs)
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			msgs <- Msg{text: s.Text(), time: time.Now()}
		}
		if err := s.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to scan input: %v\n", err)
			os.Exit(1)
		}
	}()

	finish := make(chan struct{})
	go func() {
		var initTime time.Time
		for msg := range msgs {
			log(msg, &initTime)
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
