package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2016 Essential Kaos                         //
//      Essential Kaos Open Source License <http://essentialkaos.com/ekol?en>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v5/env"
	"pkg.re/essentialkaos/ek.v5/timeutil"
	"pkg.re/essentialkaos/ek.v5/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	ARG_HELP = "h:help"
	ARG_VER  = "v:version"
)

// ////////////////////////////////////////////////////////////////////////////////// //

var cmd *exec.Cmd

// ////////////////////////////////////////////////////////////////////////////////// //

func main() {
	runtime.GOMAXPROCS(2)

	if len(os.Args) <= 2 {
		showUsage()
		return
	}

	switch os.Args[1] {
	case "-h", "--help":
		showUsage()
		return
	case "-v", "--version":
		showAbout()
		return
	}

	run(os.Args)
}

func parseArgs(args []string) (int64, string, []string) {
	maxWait := os.Args[1]
	cmdTarget := os.Args[2]
	cmdArgs := os.Args[3:]

	waitSec := timeutil.ParseDuration(maxWait)

	if waitSec == 0 {
		fmt.Fprintln(os.Stderr, "Can't parse max-time argument")
		os.Exit(99)
	}

	if env.Which(cmdTarget) == "" {
		cmdTarget, _ = filepath.Abs(args[0])
	}

	return waitSec, cmdTarget, cmdArgs
}

func run(args []string) {
	maxWait, cmdTarget, cmdArgs := parseArgs(args)

	go monitor(maxWait)

	cmd = exec.Command(cmdTarget, cmdArgs...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if cmd.Run() == nil {
		os.Exit(0)
	}

	status, ok := cmd.ProcessState.Sys().(syscall.WaitStatus)

	if ok {
		os.Exit(status.ExitStatus())
	}

	os.Exit(1)
}

func monitor(maxWait int64) {
	elapsed := int64(0)

	for {
		time.Sleep(time.Second)

		elapsed++

		if elapsed >= maxWait {
			cmd.Process.Kill()
			os.Exit(99)
		}
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "max-time", "command...")

	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.AddExample(
		"5m ./my-script.sh arg1 arg2",
		"Run my-script.sh with 5 minute limit",
	)

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     "Deadline",
		Version: "1.0.0",
		Desc:    "Simple utility for controlling application working time",
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
