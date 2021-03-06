<p align="center"><a href="#readme"><img src="https://gh.kaos.st/deadline.svg"/></a></p>

<p align="center">
  <a href="https://github.com/essentialkaos/deadline/actions"><img src="https://github.com/essentialkaos/deadline/workflows/CI/badge.svg" alt="GitHub Actions Status" /></a>
  <a href="https://github.com/essentialkaos/deadline/actions?query=workflow%3ACodeQL"><img src="https://github.com/essentialkaos/deadline/workflows/CodeQL/badge.svg" /></a>
  <a href="https://goreportcard.com/report/github.com/essentialkaos/deadline"><img src="https://goreportcard.com/badge/github.com/essentialkaos/deadline"></a>
  <a href="https://codebeat.co/projects/github-com-essentialkaos-deadline-master"><img alt="codebeat badge" src="https://codebeat.co/badges/698e5d36-2465-4266-b3d2-7f58e52d5362" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`deadline` is a simple utility for controlling application working time. Unlike [`timeout`](https://linux.die.net/man/1/timeout), `deadline` sends `KILL` signal for main processes and all child processes. This feature is very useful for shell scripts.

### Installation

#### From sources

To build the `deadline` from scratch, make sure you have a working Go 1.14+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go get github.com/essentialkaos/deadline
```

If you want to update `deadline` to latest stable release, do:

```
go get -u github.com/essentialkaos/deadline
```

#### From [ESSENTIAL KAOS Public Repository](https://yum.kaos.st)

```bash
sudo yum install -y https://yum.kaos.st/get/$(uname -r).rpm
sudo yum install deadline
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/deadline/latest).

To install the latest prebuilt version, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) deadline
```

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

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
