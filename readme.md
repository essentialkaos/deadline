<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

## `deadline`

`deadline` is a simple utility for controlling application working time.

### Installation

<details>
<summary><strong>From sources</strong></summary>

To build the MDToc from scratch, make sure you have a working Go 1.5+ workspace ([instructions](https://golang.org/doc/install)), then:

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
