package server

import (
	"time"

	"src.acicovic.me/koipond/set"
)

// Database is an item store.
//
// Database is initialized during server bootstrap, and should be
// considered as R/O after that by all threads that need to access
// the data. Concurrent R/W operations are not thread-safe.
type Database struct {
	created      time.Time
	lastModified time.Time

	items []*Item

	enabledTypes set.Strings
	defaults     map[string]string
}

// Item is a generic object stored in the database identified by
// a unique ID, with a non-empty label attached to it, and a set
// of key-value Metadata pairs.
//
// Item's Type is used to determine how various operations are
// performed on items of the same type, as well as to group items
// by type in catalogues. Generic types (any enabled type in the db
// with a valid name) are handled the same, but there are hard-coded
// specific behaviours for certain special types that change the way
// some operations are performed depending on the type,
// e.g. the SortItems function.
type Item struct {
	ID       int               `json:"id"`
	Type     string            `json:"type"`
	Label    string            `json:"label"`
	Metadata map[string]string `json:"metadata"`
}

// Catalogue is a collection of items grouped by type.
type Catalogue struct {
	Groups map[string][]*Item `json:"groups"`
}
