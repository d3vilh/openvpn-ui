package cexp

import "gopkg.in/hlandau/measurable.v1"
import "sync/atomic"

// Counter

type Counter struct {
	name  string
	value int64
}

func (c *Counter) MsName() string {
	return c.name
}

func (c *Counter) MsInt64() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *Counter) Get() int64 {
	return c.MsInt64()
}

// v must be non-negative.
func (c *Counter) Add(v int64) {
	atomic.AddInt64(&c.value, v)
}

func (c *Counter) Inc() {
	c.Add(1)
}

func (c *Counter) MsType() measurable.Type {
	return measurable.CounterType
}

func NewCounter(name string) *Counter {
	c := &Counter{
		name: name,
	}

	measurable.Register(c)
	return c
}

// Gauge

type Gauge struct {
	name  string
	value int64
}

func (c *Gauge) MsName() string {
	return c.name
}

func (c *Gauge) MsInt64() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *Gauge) Add(v int64) {
	atomic.AddInt64(&c.value, v)
}

func (c *Gauge) Sub(v int64) {
	c.Add(-v)
}

func (c *Gauge) Set(v int64) {
	atomic.StoreInt64(&c.value, v)
}

func (c *Gauge) Get() int64 {
	return c.MsInt64()
}

func (c *Gauge) Inc() {
	c.Add(1)
}

func (c *Gauge) Dec() {
	c.Add(-1)
}

func (c *Gauge) MsType() measurable.Type {
	return measurable.GaugeType
}

func NewGauge(name string) *Gauge {
	c := &Gauge{
		name: name,
	}

	measurable.Register(c)
	return c
}
