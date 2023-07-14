# Go sigar

## Overview

Go sigar is a golang implementation of the
[sigar API](https://github.com/hyperic/sigar).  The Go version of
sigar has a very similar interface, but is being written from scratch
in pure go/cgo, rather than cgo bindings for libsigar.

## Test drive

    $ git clone https://github.com/cloudfoundry/gosigar.git
    $ cd gosigar/examples
    $ go run uptime.go
    $ go run df.go
    $ go run free.go
    $ go run ps.go

## Supported platforms

Currently targeting modern flavors of macOS (Darwin), Windows and Linux.

## License

Apache 2.0
