// Package measurable provides a functionality-free integration nexus for
// metric registration.
//
// Measurable is a Go package for connecting service metrics and metric consumers.
//
// The most noteworthy feature of measurable is that it doesn't do anything.
// It contains no functionality for defining or exporting metrics.
//
// The purpose of measurable is to act as an integration nexus
// (https://www.devever.net/~hl/nexuses), essentially a matchmaker between
// application metrics and metric consumers. This creates the important feature
// that your application's metrics can be defined completely independently of
// *how* those metrics are defined.
//
// Measurable doesn't implement any metric definition or export logic because it
// strives to be a neutral intermediary, which abstracts the interface between
// measurables and measurable consumers
//
// Pursuant to this, package measurable is this and only this: an interface
// Measurable which all metrics must implement, and a facility for registering
// Measurables and visiting them.
package measurable // import "gopkg.in/hlandau/measurable.v1"

import "sync"
import "fmt"

// Measurable is the interface which must be implemented by any metric item to
// be used with package measurable. In the current version, v1, it contains
// only the MsName() and MsType() methods. All other functionality must be
// obtained by interface upgrades.
type Measurable interface {
	// Returns the name of the metric. Names should be in the style
	// "alpha.beta.gamma-delta", for example "foo.http.requests.count". That is,
	// names should be lowercase, should express a hierarchy separated by dots,
	// and have words separated by dashes.
	//
	// Some Measurable consumers may mutate these names to satisfy naming
	// restrictions applied by some graphing systems.
	MsName() string

	// Return the Measurable type. You can, of course, invent your own Measurable
	// types, though consumers won't necessarily know what to do with them.
	MsType() Type
}

var measurablesMutex sync.RWMutex
var measurables = map[string]Measurable{}

// Registers a top-level Configurable.
func Register(measurable Measurable) {
	measurablesMutex.Lock()
	defer measurablesMutex.Unlock()

	if measurable == nil {
		panic("cannot register nil measurable")
	}

	name := measurable.MsName()
	if name == "" {
		panic("measurable cannot have empty name")
	}

	_, exists := measurables[name]
	if exists {
		panic(fmt.Sprintf("A measurable with the same name already exists: %s", name))
	}

	measurables[name] = measurable
	callRegistrationHooks(measurable, RegisterEvent)
}

func Unregister(measurableName string) {
	measurablesMutex.Lock()
	defer measurablesMutex.Unlock()

	measurable, ok := measurables[measurableName]
	if !ok {
		return
	}

	callRegistrationHooks(measurable, UnregisterEvent)
	delete(measurables, measurableName)
}

func Get(measurableName string) Measurable {
	measurablesMutex.RLock()
	defer measurablesMutex.RUnlock()

	return measurables[measurableName]
}

// Visits all registered top-level Measurables.
//
// Returning a non-nil error short-circuits the iteration process and returns
// that error.
func Visit(do func(measurable Measurable) error) error {
	measurablesMutex.Lock()
	defer measurablesMutex.Unlock()

	for _, measurable := range measurables {
		err := do(measurable)
		if err != nil {
			return err
		}
	}

	return nil
}

// Represents a measurable type.
type Type uint32

const (
	// A CounterType Measurable represents a non-negative integral value
	// that monotonously increases. It must implement `MsInt64() int64`.
	CounterType Type = 0x436E7472

	// A GaugeType Measurable represents an integral value that varies over
	// time. It must implement `MsInt64() int64`.
	GaugeType = 0x47617567
)

// Registration hooks.
type HookEvent int

const (
	// This event is issued when a measurable is registered.
	RegisterEvent HookEvent = iota

	// This event is issued when a registration hook is registered. It is issued
	// for every measurable which has already been registered.
	RegisterCatchupEvent

	// This event is issued when a measurable is unregistered.
	UnregisterEvent
)

type HookFunc func(measurable Measurable, hookEvent HookEvent)

var hooksMutex sync.RWMutex
var hooks = map[interface{}]HookFunc{}

// Register for notifications on metric registration. The key must be usable as
// a key in a map and identifies the hook. No other hook with the same key must
// already exist.
//
// NOTE: The hook will be called for all registrations which already exist.
// This ensures that no registrations are missed in a threadsafe manner.
// For these calls, the event will be EventRegisterCatchup.
//
// The hook must not register or unregister registration hooks or metrics.
func RegisterHook(key interface{}, hook HookFunc) {
	measurablesMutex.RLock()
	defer measurablesMutex.RUnlock()

	registerHook(key, hook)

	for _, m := range measurables {
		hook(m, RegisterCatchupEvent)
	}
}

func registerHook(key interface{}, hook HookFunc) {
	hooksMutex.Lock()
	defer hooksMutex.Unlock()

	_, exists := hooks[key]
	if exists {
		panic(fmt.Sprintf("A metric registration hook with the same key already exists: %+v", key))
	}

	hooks[key] = hook
}

// Unregister an existing hook.
func UnregisterHook(key interface{}) {
	hooksMutex.Lock()
	defer hooksMutex.Unlock()
	delete(hooks, key)
}

func callRegistrationHooks(measurable Measurable, event HookEvent) {
	hooksMutex.RLock()
	defer hooksMutex.RUnlock()

	for _, v := range hooks {
		v(measurable, event)
	}
}
