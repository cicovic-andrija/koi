package server

import (
	"strings"
)

// RenderingCustomizer defines options available during page rendering.
type RenderingCustomizer struct {
	flags map[string]bool
}

// DataObjectInterface defines functions that need to be implemented by data types
// that are accessed by the template engine during rendering.
type DataObjectInterface interface {
	// Ref returns the reference to self or to an internal "raw" data object.
	Ref() any

	// Properties returns the reference to a map of object's properties.
	Properties() map[string]string

	// Groups returns groups of other data objects, for container data types.
	Groups() map[string][]DataObjectInterface

	// MultiGroup returns a flag indicating whether there is more than one group.
	MultiGroup() bool

	// Tags returns object's tags.
	Tags() []string

	// HideTags returns a flag indicating whether tags should not be rendered.
	HideTags() bool
}

// CommonBaseObject provides a dummy implementation of DataObjectInterface.
type CommonBaseObject struct{}

// Test returns the boolean value of flag.
func (r *RenderingCustomizer) Test(flag string) bool {
	return r.flags[flag]
}

// Capitalize returns the passed string after converting it to uppercase.
func (r *RenderingCustomizer) Capitalize(s string) string {
	return strings.ToUpper(s)
}

// Ref implements DataObjectInterface.
func (*CommonBaseObject) Ref() any {
	return nil
}

// Properties implements DataObjectInterface.
func (*CommonBaseObject) Properties() map[string]string {
	return nil
}

// Groups implements DataObjectInterface.
func (*CommonBaseObject) Groups() map[string][]DataObjectInterface {
	return nil
}

// MultiGroup implements DataObjectInterface.
func (*CommonBaseObject) MultiGroup() bool {
	return false
}

// Tags implements DataObjectInterface.
func (*CommonBaseObject) Tags() []string {
	return nil
}

// HideTags implements DataObjectInterface.
func (*CommonBaseObject) HideTags() bool {
	return false
}
