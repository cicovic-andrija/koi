package server

import (
	"time"

	"src.acicovic.me/koipond/set"
)

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

// Database is an item store.
//
// Database is initialized during server bootstrap, and should be
// considered as R/O after that by all threads that need to access
// the data. Concurrent R/W operations are not thread-safe. In fact,
// none of the Database methods are thread-safe.
type Database struct {
	filePath     string
	created      time.Time
	lastModified time.Time

	items []*Item

	enabledTypes set.Strings
	defaults     map[string]string
}

// Global database instance.
var _database = &Database{
	items:        []*Item{},
	defaults:     map[string]string{},
	enabledTypes: set.NewStringSet(),
}

// AddItem creates and adds a new item to the Database. The function
// returns a pointer to the item, or nil if the item could not be added.
func (db *Database) AddItem(typeKey string, metadata map[string]string) *Item {
	item := &Item{
		ID:       len(db.items),
		Type:     typeKey,
		Metadata: metadata,
	}
	if ok := item.SetLabel(); !ok {
		return nil
	}
	db.items = append(db.items, item)
	return item
}

// MetadataValue returns metadata associated with the given key.
// Zero value, when returned, indicates that there was no
// metadata associated with the key.
func (i *Item) MetadataValue(key string) string {
	return i.Metadata[key]
}

// SetLabel sets the Item's Label property as defined
// by its type and metadata.
func (i *Item) SetLabel() bool {
	switch i.Type {
	case "books":
		i.Label = i.Metadata["title"]
	case "games":
		i.Label = i.Metadata["title"]
	case "workouts":
		i.Label = i.Metadata["date"]
	default:
		i.Label = i.Metadata["label"]
	}
	return i.Label != ""
}

// SortItems sorts the passed slice as defined by item type.
// The function assumes that all items are of the same type.
func SortItems(items []*Item) {
	if len(items) == 0 {
		return
	}

	switch items[0].Type {
	case "books":
		By(bookTitleStandard).Sort(items)
	case "games":
		By(groupedUnderSeries).Sort(items)
	case "equipment":
		By(label).Sort(items)
	default:
		By(label).Sort(items)
	}
}

// IsValidItemKeyWordForType checks if the key is a valid single
// item keyword in the database (e.g. "book" for "books" type).
// Keyword "item" is always valid for every type.
func IsValidItemKeyWordForType(key string, typeKey string) bool {
	if key == "item" {
		return true
	}
	switch typeKey {
	case "books":
		return key == "book"
	case "games":
		return key == "game"
	case "workouts":
		return key == "workout"
	}
	return false
}
