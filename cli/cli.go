package cli

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
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

	"github.com/essentialkaos/ek/v12/env"
	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/strutil"
	"github.com/essentialkaos/ek/v12/support"
	"github.com/essentialkaos/ek/v12/support/deps"
	"github.com/essentialkaos/ek/v12/system/process"
	"github.com/essentialkaos/ek/v12/terminal"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/timeutil"
	"github.com/essentialkaos/ek/v12/usage"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	APP  = "deadline"
	VER  = "1.6.2"
	DESC = "Simple utility for controlling application working time"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	OPT_HELP = "h:help"
	OPT_VER  = "v:version"
)

// ERROR_EXIT_CODE is exit code used as error marker
const ERROR_EXIT_CODE = 255

// ////////////////////////////////////////////////////////////////////////////////// //

type SignalInfo struct {
	Wait   time.Duration
	Signal syscall.Signal
}

// ////////////////////////////////////////////////////////////////////////////////// //

// cmd is passed command
var cmd *exec.Cmd

// colorTagApp contains color tag for app name
var colorTagApp string

// colorTagVer contains color tag for app version
var colorTagVer string

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main application function
func Run(gitRev string, gomod []byte) {
	runtime.GOMAXPROCS(2)

	preConfigureUI()

	if len(os.Args) == 2 {
		switch os.Args[1] {
		case "-h", "--help":
			genUsage().Print()
			return
		case "-v", "--version":
			genAbout(gitRev).Print()
			return
		case "-vv", "--verbose-version":
			support.Collect(APP, VER).
				WithRevision(gitRev).
				WithDeps(deps.Extract(gomod)).
				Print()
			return
		}
	}

	if len(os.Args) <= 2 {
		genUsage().Print()
		return
	}

	run(os.Args)
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}

	switch {
	case fmtc.Is256ColorsSupported():
		colorTagApp, colorTagVer = "{*}{#160}", "{#160}"
	default:
		colorTagApp, colorTagVer = "{*}{r}", "{r}"
	}
}

// run run application and monitor
func run(args []string) {
	sigInfo, cmdApp, cmdArgs := parseArgs(args)

	go monitor(sigInfo)

	cmd = exec.Command(cmdApp, cmdArgs...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

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
		terminal.Error(err)
		os.Exit(ERROR_EXIT_CODE)
	}

	cmdApp = env.Which(args[2])

	if cmdApp == "" {
		cmdApp, _ = filepath.Abs(args[2])
	}

	if !fsutil.CheckPerms("FX", cmdApp) {
		terminal.Error("Can't execute command %q", args[2])
		os.Exit(ERROR_EXIT_CODE)
	}

	if len(args) > 3 {
		cmdArgs = args[3:]
	}

	return sigInfo, cmdApp, cmdArgs
}

// monitor send signal if elapsed time greater than max time
func monitor(sigInfo SignalInfo) {
	var signalSent bool

	start := time.Now()

	for range time.NewTicker(time.Second).C {
		if time.Since(start) >= sigInfo.Wait && !signalSent {
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

	if len(info.Children) == 0 {
		return result
	}

	for _, p := range info.Children {
		result = append(result, getAllSubProcPIDs(p)...)
	}

	return result
}

// parseSignalInfo parse signal data
func parseTimeAndSignal(data string) (SignalInfo, error) {
	var err error
	var wait time.Duration
	var sig syscall.Signal

	if !strings.Contains(data, ":") {
		wait, err = timeutil.ParseDuration(data)

		if err != nil {
			return SignalInfo{}, fmt.Errorf("Can't parse %q", data)
		}

		return SignalInfo{wait, syscall.SIGTERM}, nil
	}

	wait, err = timeutil.ParseDuration(strutil.ReadField(data, 0, true, ':'))

	if err != nil {
		return SignalInfo{}, fmt.Errorf("Can't parse %q", data)
	}

	signal := strutil.ReadField(data, 1, true, ':')

	switch strings.ToUpper(signal) {
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
		return SignalInfo{}, fmt.Errorf("Unsupported signal %s", signal)
	}

	return SignalInfo{wait, sig}, nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "time:signal", "command…")

	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"5m my-script.sh arg1 arg2",
		"Run my-script.sh and send TERM signal in 5 minutes",
	)

	info.AddExample(
		"5m:KILL my-script.sh arg1 arg2",
		"Run my-script.sh and send KILL signal in 5 minutes",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,

		AppNameColorTag: colorTagApp,
		VersionColorTag: colorTagVer,
		DescSeparator:   "{s}—{!}",

		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
