<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

## `deadline` [![Build Status](https://travis-ci.org/essentialkaos/deadline.svg?branch=master)](https://travis-ci.org/essentialkaos/deadline) [![Go Report Card](https://goreportcard.com/badge/github.com/essentialkaos/deadline)](https://goreportcard.com/report/github.com/essentialkaos/deadline) [![License](https://gh.kaos.io/ekol.svg)](https://essentialkaos.com/ekol)

`deadline` is a simple utility for controlling application working time. Unlike [`timeout`](https://linux.die.net/man/1/timeout) `deadline` send `KILL` signal for main processes and all child processes. This feature is very useful with shell scripts.

### Installation

<details>
<summary><strong>From sources</strong></summary>

To build the `deadline` from scratch, make sure you have a working Go 1.5+ workspace ([instructions](https://golang.org/doc/install)), then:

```
go get github.com/essentialkaos/deadline
```

If you want update `deadline` to latest stable release, do:

```
go get -u github.com/essentialkaos/deadline
```
</details>

<details>
<summary><strong>From ESSENTIAL KAOS Public repo for RHEL6/CentOS6</strong></summary>
```
[sudo] yum install -y https://yum.kaos.io/6/release/i386/kaos-repo-7.2-0.el6.noarch.rpm
[sudo] yum install deadline
```
</details>

<details>
<summary><strong>From ESSENTIAL KAOS Public repo for RHEL7/CentOS7</strong></summary>
```
[sudo] yum install -y https://yum.kaos.io/7/release/x86_64/kaos-repo-7.2-0.el7.noarch.rpm
[sudo] yum install deadline
```
</details>

### Usage

```
Usage: deadline {options} time:signal command...

Options

  --help, -h       Show this help message
  --version, -v    Show version

Examples

  deadline 5m my-script.sh arg1 arg2
  Run my-script.sh and send TERM signal in 5 minutes

  deadline 5m:KILL my-script.sh arg1 arg2
  Run my-script.sh and send KILL signal in 5 minutes

```

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)
