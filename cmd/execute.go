package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

func dummyExecute(bin string, args ...string) error {
	cmd := exec.Command(bin, args...)
	fmt.Println(cmd.String())
	return nil
}

func execute(bin string, args ...string) error {
	var wg sync.WaitGroup

	cmd := exec.Command(bin, args...)
	cmd.String()

	// pipe StdOut to stdout
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	go logLine(stdOut)

	// pipe StdErr to stdout
	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go logLine(stdErr)

	// start process
	err = cmd.Start()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}

	wg.Add(2)

	go func() {
		defer wg.Done()
		logLine(stdOut)
	}()

	go func() {
		defer wg.Done()
		logLine(stdErr)
	}()

	wg.Wait()

	// Wait for the command to finish
	return cmd.Wait()
}

func logLine(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			line := string(buf[0:n])
			fmt.Println(line)
		}
		if err != nil {
			break
		}
	}
}
