package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// go run main.go run 		<command> <args>
// docker		run <image> <command> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("I am confused")
	}

}

func run() {
	fmt.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER, //separate hostname from namespace hostname | creating usernamespace for child process
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      1000,
				Size:        1,
			},
		},
	}
	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
