<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/deadline/ci"><img src="https://kaos.sh/w/deadline/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/deadline"><img src="https://kaos.sh/r/deadline.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/deadline"><img src="https://kaos.sh/b/698e5d36-2465-4266-b3d2-7f58e52d5362.svg" alt="codebeat badge" /></a>
  <a href="https://kaos.sh/w/deadline/codeql"><img src="https://kaos.sh/w/deadline/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`deadline` is a simple utility for controlling application working time. Unlike [`timeout`](https://linux.die.net/man/1/timeout), `deadline` sends `KILL` signal for main processes and all child processes. This feature is very useful for shell scripts.

### Installation

#### From sources

To build the `deadline` from scratch, make sure you have a working Go 1.21+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/deadline@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://kaos.sh/kaos-repo)

```bash
sudo yum install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install deadline
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux from [EK Apps Repository](https://apps.kaos.st/deadline/latest).

To install the latest prebuilt version, do:

```bash
bash <(curl -fsSL https://apps.kaos.st/get) deadline
```

### Usage

<p align="center"><img src=".github/images/usage.svg"/></p>

### Build Status

| Branch | Status |
|--------|--------|
| `master` | [![CI](https://kaos.sh/w/deadline/ci.svg?branch=master)](https://kaos.sh/w/deadline/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/deadline/ci.svg?branch=master)](https://kaos.sh/w/deadline/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
