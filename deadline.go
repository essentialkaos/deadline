// +build linux

package main

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"pkg.re/essentialkaos/ek.v8/env"
	"pkg.re/essentialkaos/ek.v8/fsutil"
	"pkg.re/essentialkaos/ek.v8/system/process"
	"pkg.re/essentialkaos/ek.v8/timeutil"
	"pkg.re/essentialkaos/ek.v8/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "deadline"
	VER  = "1.3.1"
	DESC = "Simple utility for controlling application working time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	ARG_HELP = "h:help"
	ARG_VER  = "v:version"
)

// exit code used as error marker
const ERROR_EXIT_CODE = 99

// ////////////////////////////////////////////////////////////////////////////////// //

type SignalInfo struct {
	Wait   int64
	Signal syscall.Signal
}

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

// run run application and monitor
func run(args []string) {
	sigInfo, cmdApp, cmdArgs := parseArgs(args)

	go monitor(sigInfo)

	cmd = exec.Command(cmdApp, cmdArgs...)

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

	os.Exit(ERROR_EXIT_CODE)
}

// parseArgs parse command-line arguments
func parseArgs(args []string) (SignalInfo, string, []string) {
	var err error
	var sigInfo SignalInfo
	var cmdApp string
	var cmdArgs []string

	sigInfo, err = parseTimeAndSignal(args[1])

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(ERROR_EXIT_CODE)
	}

	cmdApp = env.Which(args[2])

	if cmdApp == "" {
		cmdApp, _ = filepath.Abs(args[2])
	}

	if !fsutil.CheckPerms("FX", cmdApp) {
		fmt.Fprintf(os.Stderr, "Can't execute %s\n", args[2])
		os.Exit(ERROR_EXIT_CODE)
	}

	if len(args) > 3 {
		cmdArgs = args[3:]
	}

	return sigInfo, cmdApp, cmdArgs
}

// monitor send signal if elapsed time greater than max time
func monitor(sigInfo SignalInfo) {
	var elapsed int64
	var signalSent bool

	for {
		time.Sleep(time.Second)

		elapsed++

		if elapsed >= sigInfo.Wait && !signalSent {
			signalSent = true
			sendSignal(sigInfo.Signal)
		}
	}
}

// sendSignal send signal to application
func sendSignal(signal syscall.Signal) {
	if signal == syscall.SIGKILL {
		killAllProcesses()
		os.Exit(ERROR_EXIT_CODE)
	} else {
		syscall.Kill(cmd.Process.Pid, signal)
	}
}

// killAllProcesses kill application with all subprocesses
func killAllProcesses() {
	tree, err := process.GetTree(cmd.Process.Pid)

	if err != nil {
		syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)
		os.Exit(ERROR_EXIT_CODE)
	}

	list := getAllSubProcPIDs(tree)

	syscall.Kill(cmd.Process.Pid, syscall.SIGKILL)

	for _, pid := range list {
		pidDir := fmt.Sprintf("/proc/%d", pid)

		if fsutil.IsExist(pidDir) {
			syscall.Kill(pid, syscall.SIGKILL)
		}
	}

	os.Exit(ERROR_EXIT_CODE)
}

// getAllSubProcPIDs return slice with all pids in process tree
func getAllSubProcPIDs(info *process.ProcessInfo) []int {
	var result = []int{info.PID}

	if len(info.Childs) == 0 {
		return result
	}

	for _, p := range info.Childs {
		result = append(result, getAllSubProcPIDs(p)...)
	}

	return result
}

// parseSignalInfo parse signal data
func parseTimeAndSignal(data string) (SignalInfo, error) {
	var wait int64
	var sig syscall.Signal

	if !strings.Contains(data, ":") {
		wait = timeutil.ParseDuration(data)

		if wait == 0 {
			return SignalInfo{}, fmt.Errorf("Can't parse %s", data)
		}

		return SignalInfo{wait, syscall.SIGTERM}, nil
	}

	dataSlice := strings.Split(data, ":")

	wait = timeutil.ParseDuration(dataSlice[0])

	if wait == 0 {
		return SignalInfo{}, fmt.Errorf("Can't parse %s", data)
	}

	switch strings.ToUpper(dataSlice[1]) {
	case "SIGABRT", "ABRT", "6":
		sig = syscall.SIGABRT
	case "SIGALRM", "ALRM", "14":
		sig = syscall.SIGALRM
	case "SIGBUS", "BUS", "10":
		sig = syscall.SIGBUS
	case "SIGCHLD", "CHLD", "18":
		sig = syscall.SIGCHLD
	case "SIGCONT", "CONT", "25":
		sig = syscall.SIGCONT
	case "SIGFPE", "FPE", "8":
		sig = syscall.SIGFPE
	case "SIGHUP", "HUP", "1":
		sig = syscall.SIGHUP
	case "SIGILL", "ILL", "4":
		sig = syscall.SIGILL
	case "SIGINT", "INT", "2":
		sig = syscall.SIGINT
	case "SIGIO", "IO":
		sig = syscall.SIGIO
	case "SIGIOT", "IOT":
		sig = syscall.SIGIOT
	case "SIGKILL", "KILL", "9":
		sig = syscall.SIGKILL
	case "SIGPIPE", "PIPE", "13":
		sig = syscall.SIGPIPE
	case "SIGPROF", "PROF", "29":
		sig = syscall.SIGPROF
	case "SIGQUIT", "QUIT", "3":
		sig = syscall.SIGQUIT
	case "SIGSEGV", "SEGV", "11":
		sig = syscall.SIGSEGV
	case "SIGSTOP", "STOP", "23":
		sig = syscall.SIGSTOP
	case "SIGSYS", "SYS", "12":
		sig = syscall.SIGSYS
	case "SIGTERM", "TERM", "15":
		sig = syscall.SIGTERM
	case "SIGTRAP", "TRAP", "5":
		sig = syscall.SIGTRAP
	case "SIGTSTP", "TSTP", "20":
		sig = syscall.SIGTSTP
	case "SIGTTIN", "TTIN", "26":
		sig = syscall.SIGTTIN
	case "SIGTTOU", "TTOU", "27":
		sig = syscall.SIGTTOU
	case "SIGURG", "URG", "21":
		sig = syscall.SIGURG
	case "SIGUSR1", "USR1", "16":
		sig = syscall.SIGUSR1
	case "SIGUSR2", "USR2", "17":
		sig = syscall.SIGUSR2
	case "SIGVTALRM", "VTALRM", "28":
		sig = syscall.SIGVTALRM
	case "SIGWINCH", "WINCH":
		sig = syscall.SIGWINCH
	case "SIGXCPU", "XCPU", "30":
		sig = syscall.SIGXCPU
	case "SIGXFSZ", "XFSZ", "31":
		sig = syscall.SIGXFSZ
	case "":
		return SignalInfo{}, fmt.Errorf("Signal is not set")
	default:
		return SignalInfo{}, fmt.Errorf("Unsupported signal %s", dataSlice[1])
	}

	return SignalInfo{wait, sig}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo("", "time:signal", "command...")

	info.AddOption(ARG_HELP, "Show this help message")
	info.AddOption(ARG_VER, "Show version")

	info.AddExample(
		"5m my-script.sh arg1 arg2",
		"Run my-script.sh and send TERM signal in 5 minutes",
	)

	info.AddExample(
		"5m:KILL my-script.sh arg1 arg2",
		"Run my-script.sh and send KILL signal in 5 minutes",
	)

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
