package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

func main() {
	name := flag.String("name", "", "Name of the command to execute")
	uid := flag.String("uid", "0", "User ID to run the command as")
	gid := flag.String("gid", "3005", "Group ID to run the command as")
	outFile := flag.String("out", "null", "File to write stdout to")
	errFile := flag.String("err", "null", "File to write stderr to")
	flag.Parse()
	args := flag.Args()
	if *name == "" {
		fmt.Println("Error: -name parameter is required")
		os.Exit(1)
	}
	os.Exit(startproc(*name, *uid, *gid, *outFile, *errFile, args...))
}

func startproc(name string, uid string, gid string, outFile string, errFile string, arg ...string) (exitCode int) {
	var out, err *os.File
	var openErr error
	if outFile == "null" {
		out = nil
	} else {
		out, openErr = os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0600)
		if openErr != nil {
			fmt.Println("Error opening output file:", openErr)
			return 1
		}
	}
	if errFile == "null" {
		err = nil
	} else {
		err, openErr = os.OpenFile(errFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_TRUNC, 0600)
		if openErr != nil {
			fmt.Println("Error opening error file:", openErr)
			return 1
		}
	}

	proc := exec.Command(name, arg...)
	proc.Stdout = out
	proc.Stderr = err

	proc.SysProcAttr = &syscall.SysProcAttr{}
	if proc.SysProcAttr == nil {
		fmt.Println("Failed to initialize SysProcAttr")
	}

	proc.SysProcAttr.Setpgid = true
	if proc.SysProcAttr == nil {
		fmt.Println("Failed to set process group ID")
	}

	uidInt, erro := strconv.Atoi(uid)
	if erro != nil {
		fmt.Println("Error converting uid to integer:", erro)
		return 1
	}
	gidInt, erro := strconv.Atoi(gid)
	if erro != nil {
		fmt.Println("Error converting gid to integer:", erro)
		return 1
	}

	proc.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uidInt), Gid: uint32(gidInt)}
	if proc.SysProcAttr.Credential == nil {
		fmt.Println("Failed to set credentials")
	}

	if erro := proc.Start(); erro != nil {
		fmt.Fprintf(os.Stderr, "Error starting process: %v\n", erro)
		return 1
	}
	return 0
}
