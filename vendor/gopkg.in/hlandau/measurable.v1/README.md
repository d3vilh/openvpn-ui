Measurable: The useless Go metric registration package that doesn't do anything
===============================================================================

[![GoDoc](https://godoc.org/gopkg.in/hlandau/measurable.v1?status.svg)](https://godoc.org/gopkg.in/hlandau/measurable.v1)

Measurable is a Go library for managing the registration of metrics such as
counters and gauges, no matter how that metric data is eventually consumed.

The most noteworthy feature of measurable is that it doesn't do anything. It
contains no functionality for providing metric data to any external service,
and it contains no actual metric implementations.

The purpose of measurable is to act as an [integration
nexus](https://www.devever.net/~hl/nexuses), essentially a matchmaker between
metric sources and metric consumers. This creates the important feature that
your application's metrics can be expressed completely independently of *how*
those metrics are exported.

Measurable doesn't implement any metric or metric export logic because it
strives to be a neutral intermediary, which abstracts the interface between
metrics and metric exporters.

**Import as:** `gopkg.in/hlandau/measurable.v1`

Measurable
----------

A Measurable is an object that represents some metric. It is obliged only to
implement the following interface:

```go
type Measurable interface {
  MsName() string
  MsType() Type
}
```

Measurable is designed around interface upgrades. If you want to actually
do anything with a Measurable, you must attempt to cast it to an interface
with the methods you need. A Measurable is not obliged to implement any
interface besides Measurable, but almost always will.

Here are some common interfaces implemented by Measurables, in descending order
of importance:

  - `MsName() string` — get the Measurable name.
  - `MsType() Type` — get the Measurable type.
  - `MsInt64() int64` — get the Measurable as an int64.
  - `String() string` — the standard Go `String()` interface.

All Measurables should implement `MsName() string` and `MsType() Type`.

Measurable-specific methods should always be prefixed by `Ms` so it is clear
they are intended for consumption by Measurable consumers.

`MsName`, `MsType` and `MsInt64` should suffice for most consumers of Counter
and Gauge metric types.

Metrics should be named in lowercase using dots to create a hierarchy and
dashes to separate words, e.g. `someserver.http.request-count`. These metric
names may be transmuted by consumers as necessary for some graphing systems,
such as Prometheus (which allows only underscores).

Standard Bindings
-----------------

For a package which makes it easy to register and consume measurables, see the
[easymetric](https://github.com/hlandau/easymetric) package.

Of course, nothing requires you to use the easymetric package. You are free to escew it and make your own.

Background Reading
------------------

  - [On Nexuses](https://www.devever.net/~hl/nexuses)
  - See also: [Configurable](https://github.com/hlandau/configurable)

Licence
-------

    © 2015 Hugo Landau <hlandau@devever.net>    MIT License

