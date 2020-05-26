package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

// TODO:
// 1) Try using exec.Command() instead of ForkExec()
//    More info here: https://gobyexample.com/spawning-processes
// 2) Fix output of command running right after prompt

func main() {
	// Initialize reader to read input from stdin (keyboard).
	reader := bufio.NewReader(os.Stdin)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	//exitChan := make(chan int)
	go func() {
		for _ = range signalChan {
			//fmt.Printf("\nSIGINT received: %v\n", sig)
		}
		//	signal := <-signalChan
		//	switch signal {
		//	case syscall.SIGINT:
		//		fmt.Println("SIGINT received!")
		//		exitChan <- 1
		//	default:
		//		fmt.Println("Different signal received...")
		//		exitChan <- 0
		//	}
	}()

	for {
		// Provide a prompt indicator to user.
		fmt.Print("> ")
		// Read input delimited by newline characters.
		input, err := reader.ReadString('\n')
		if err != nil {
			// Handle Ctrl+D to stop program.
			if err == io.EOF {
				fmt.Print("\nBye for now!\n")
				os.Exit(1)
			}

			fmt.Fprintln(os.Stderr, err)
		}

		//	code := <-exitChan
		//	os.Exit(code)

		input = strings.Trim(input, "\n")
		inputSlice := strings.Split(input, " ")
		//fmt.Printf("input: %v\n", inputSlice)

		switch inputSlice[0] {
		case "exit":
			os.Exit(0)
		case "sleep":
			// Use os/exec pkg here instead (and maybe throughout)
			// More info: https://stackoverflow.com/questions/11886531/terminating-a-process-started-with-os-exec-in-golang
			sleepValue := inputSlice[1]
			cmd := exec.Command("sleep", sleepValue)
			//binary, err := exec.LookPath("sleep")
			//if err != nil {
			//	fmt.Print("\nFailed to look up command\n")
			//	break
			//}

			//args := syscall.ProcAttr{
			//	"",
			//	[]string{},
			//	[]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
			//	nil,
			//}

			//_, err = syscall.ForkExec(binary, []string{"sleep", sleepValue}, &args)
			err := cmd.Run()
			var exerr *exec.ExitError
			if errors.As(err, &exerr) {
				fmt.Println()
				continue
			}
			if err != nil {
				fmt.Print("\nFailed to execute command\n")
			}
		case "pwd":
			path, err := os.Getwd()
			if err != nil {
				fmt.Printf("\nFailed to get path to current directory\n")
			}
			fmt.Println(path)
		case "ls":
			binary, err := exec.LookPath("ls")
			if err != nil {
				fmt.Print("\nFailed to look up command\n")
				break
			}

			args := syscall.ProcAttr{
				"",
				[]string{},
				[]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
				nil,
			}

			_, err = syscall.ForkExec(binary, []string{"ls"}, &args)
			if err != nil {
				fmt.Print("\nFailed to execute command\n")
			}
		case "whoami":
			binary, err := exec.LookPath("whoami")
			if err != nil {
				fmt.Print("\nFailed to look up command\n")
				break
			}

			args := syscall.ProcAttr{
				"",
				[]string{},
				[]uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
				nil,
			}

			_, err = syscall.ForkExec(binary, []string{"ls"}, &args)
			//args := []string{"whoami"}
			//env := os.Environ()

			//err = syscall.Exec(binary, args, env)
			if err != nil {
				fmt.Print("\nFailed to execute command\n")
			}
		default:
			// Echo input back to user.
			fmt.Printf("%s: command not found\n", strings.Join(inputSlice, " "))
		}
	}
}
