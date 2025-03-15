package server

import (
	"strings"
)

// RenderingCustomizer defines options available during page rendering.
type RenderingCustomizer struct {
	flags map[string]bool
}

// CommonDataProperties defines functions that need to be implemented by data types
// that are accessed by the template engine during rendering.
type CommonDataProperties interface {
	Ref() any
	Properties() map[string]string
	Groups() map[string][]CommonDataProperties
	MultiGroup() bool
	Tags() []string
	HideTags() bool
}

// CommonBaseObject provides a dummy implementation of CommonDataProperties.
type CommonBaseObject struct{}

// Test returns the boolean value of flag.
func (r *RenderingCustomizer) Test(flag string) bool {
	return r.flags[flag]
}

// Capitalize returns the passed string after converting it to uppercase.
func (r *RenderingCustomizer) Capitalize(s string) string {
	return strings.ToUpper(s)
}

// Ref implements CommonDataProperties.
func (*CommonBaseObject) Ref() any {
	return nil
}

// Properties implements CommonDataProperties.
func (*CommonBaseObject) Properties() map[string]string {
	return nil
}

// Groups implements CommonDataProperties.
func (*CommonBaseObject) Groups() map[string][]CommonDataProperties {
	return nil
}

// MultiGroup implements CommonDataProperties.
func (*CommonBaseObject) MultiGroup() bool {
	return false
}

// Tags implements CommonDataProperties.
func (*CommonBaseObject) Tags() []string {
	return nil
}

// HideTags implements CommonDataProperties.
func (*CommonBaseObject) HideTags() bool {
	return false
}
