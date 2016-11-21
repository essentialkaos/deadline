## `deadline`

`deadline` is a simple utility for controlling application working time.

### Installation

#### From sources

```
go get github.com/essentialkaos/deadline
```

#### From ESSENTIAL KAOS Public repo for RHEL6/CentOS6

```
[sudo] yum install -y https://yum.kaos.io/6/release/i386/kaos-repo-7.2-0.el6.noarch.rpm
[sudo] yum install deadline
```

#### From ESSENTIAL KAOS Public repo for RHEL7/CentOS7

```
[sudo] yum install -y https://yum.kaos.io/7/release/x86_64/kaos-repo-7.2-0.el7.noarch.rpm
[sudo] yum install deadline
```

### Usage

```
Usage: deadline max-time command...

Options

  --help, -h        Show this help message
  --version, -v     Show version

Examples

  deadline 5m ./my-script.sh arg1 arg2
  Run my-script.sh with 5 minute limit

```

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[EKOL](https://essentialkaos.com/ekol)
