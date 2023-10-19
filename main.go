package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"

	"kernel.org/pub/linux/libs/security/libcap/cap"
)

// go run main.go run 		<command> <args>
// docker		run <image> <command> <args>
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		ensureNotEUID()

		child("/mnt/app.js")
	default:
		panic("I am confused")
	}

}

func run() {

	//mount workspace
	fmt.Printf("I have just began")
	must(syscall.Mount("/home/nora/Bureau/Perso/container-from-scratch", "/home/nora/Bureau/Perso/os/ubuntu-base-14.04-core-amd64/mnt", "", syscall.MS_BIND, ""))

	fmt.Printf("\nWorkspace mounted....\n")
	//create a container

	createContainer()

	//run app in the container

}
func ensureNotEUID() {
	euid := syscall.Geteuid()
	uid := syscall.Getuid()
	egid := syscall.Getegid()
	gid := syscall.Getgid()
	if uid != euid || gid != egid {
		log.Fatalf("go runtime is setuid uids:(%d vs %d), gids(%d vs %d)", uid, euid, gid, egid)
	}
	if uid == 0 {
		log.Fatalf("go runtime is running as root - cheating")
	}

}

func child(appPath string) {

	fmt.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())

	fmt.Println("child running")
	// Get the current capabilities
	c := cap.GetProc()

	if err := c.SetFlag(cap.Effective, true, cap.SYS_CHROOT); err != nil {
		log.Fatalf("Failed to set effetive capability: %v", err)
	}
	if err := c.SetFlag(cap.Permitted, true, cap.SYS_CHROOT); err != nil {
		log.Fatalf("Failed to set permitted capability: %v", err)
	}

	// Re-check the capabilities (SYS_CHROOT should now be effective)
	c = cap.GetProc()

	log.Printf("this process has these caps: %s", c)

	// Check if the capability is granted
	if on, _ := c.GetFlag(cap.Effective, cap.SYS_CHROOT); !on {
		log.Fatalf("Insufficient effective privilege to execute syscall.Chroot - required capability not granted")
	}
	if on, _ := c.GetFlag(cap.Permitted, cap.SYS_CHROOT); !on {
		log.Fatalf("Insufficient permitted privilege to execute syscall.Chroot - required capability not granted")
	}

	// Execute the syscall.Chroot operation
	must(syscall.Chroot("/home/nora/Bureau/Perso/os/ubuntu-base-14.04-core-amd64"))

	// Remove SYS_CHROOT capability
	// if err := c.SetFlag(cap.Effective, false, cap.SYS_CHROOT); err != nil {
	// 	log.Fatalf("Failed to remove capability: %v", err)
	// }

	// Check the capabilities after removing SYS_CHROOT
	c = cap.GetProc()
	log.Printf("this process has these caps: %s", c)
	fmt.Println("child still running")
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	//mount workspace to the container

	cmd := exec.Command("/usr/local/lib/node-v10.24.1/bin/node", appPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
}

func createContainer() {
	fmt.Printf("Running %v as user %d in process %d\n", os.Args[2:], os.Geteuid(), os.Getpid())
	fmt.Println("parent calling child")
	c := cap.GetProc()

	log.Printf("this parent process has these caps: %s", c)
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
	fmt.Println("parent setting environment for child")

	must(cmd.Run())

	fmt.Println("parent done")

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
