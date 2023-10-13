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
	case "child":
		child("app.js")
	default:
		panic("I am confused")
	}

}

func run() {

	//create a container
	createContainer()

	//run app in the container

}

func child(appPath string) {
	fmt.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())

	must(syscall.Chroot("/home/nora/Bureau/Perso/os/ubuntu-base-14.04-core-amd64"))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command("/usr/local/node/node-v10.24.1-linux-x64/bin/node", appPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
}

func createContainer() {
	fmt.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID, //separate hostname from namespace hostname | creating usernamespace for child process
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
