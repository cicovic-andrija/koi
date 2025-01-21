package server

import (
	"strings"
	"time"

	"src.acicovic.me/koipond/set"
)

// Special metadata keys.
const (
	MDLabelKey       = "label"
	MDCollectionsKey = "collections"
	MDTagsKey        = "tags"
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
// e.g. the Sort function.
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

	items        []*Item
	collectioned map[string][]*Item
	tagged       map[string][]*Item

	enabledTypes        set.Strings
	declaredCollections map[string]string
	hiddenCollections   set.Strings
	defaults            map[string]string
}

// Global database instance.
var _database = &Database{
	items:               []*Item{},
	collectioned:        map[string][]*Item{},
	tagged:              map[string][]*Item{},
	enabledTypes:        set.NewStringSet(),
	declaredCollections: map[string]string{},
	hiddenCollections:   set.NewStringSet(),
	defaults:            map[string]string{},
}

// Creates and adds a new item to the Database. The function
// returns a pointer to the item, or nil if the item could not be added.
func (db *Database) add(typeKey string, metadata map[string]string) *Item {
	item := &Item{
		ID:       len(db.items),
		Type:     typeKey,
		Metadata: metadata,
	}

	if ok := item.setLabel(); !ok {
		return nil
	}

	// set default metadata where it's missing
	for key, defaultValue := range db.defaults {
		if strings.HasPrefix(key, typeKey+"/") {
			key = strings.TrimPrefix(key, typeKey+"/")
			if item.Metadata[key] == "" {
				item.Metadata[key] = defaultValue
			}
		}
	}

	// item is valid at this point, anything optional goes below
	db.items = append(db.items, item)

	// index collection (invalid collection names are ignored)
	if collections := metadata[MDCollectionsKey]; collections != "" {
		itemsCollections := []string{}
		for _, collectionKey := range strings.Split(collections, ",") {
			collectionKey = strings.TrimSpace(collectionKey)
			if isValidCollectionKey(collectionKey) {
				_, declared := db.declaredCollections[collectionKey]
				if declared && !db.hiddenCollections.Contains(collectionKey) {
					itemsCollections = append(itemsCollections, collectionKey)
				}
			}
		}
		for _, c := range itemsCollections {
			db.collectioned[c] = append(db.collectioned[c], item)
		}
	}

	// index tags (invalid tag names are ignored)
	if tags := metadata[MDTagsKey]; tags != "" {
		itemsTags := []string{}
		for _, t := range strings.Split(tags, ",") {
			t := strings.TrimSpace(t)
			if isValidTag(t) {
				itemsTags = append(itemsTags, t)
			}
		}
		for _, t := range itemsTags {
			db.tagged[t] = append(db.tagged[t], item)
		}
	}

	return item
}

func (db *Database) collections() map[string]string {
	collections := make(map[string]string)
	for key := range db.collectioned {
		collections[key] = db.declaredCollections[key]
	}
	return collections
}

func (db *Database) tags() map[string]int {
	tags := map[string]int{}
	for tag, items := range db.tagged {
		tags[tag] = len(items)
	}
	return tags
}

func (db *Database) lastID() int {
	return len(db.items) - 1
}

func (db *Database) singleItem(id int) *Item {
	return db.items[id]
}

func (db *Database) collectionCatalogue(key string) *Catalogue {
	return MakeCatalogue(db.collectioned[key])
}

func (db *Database) catalogueOfTaggedItems(tag string) *Catalogue {
	return MakeCatalogue(db.tagged[tag])
}

// Tags returns a slice of item's tags.
func (i *Item) Tags() (tags []string) {
	for _, tag := range strings.Split(i.Metadata[MDTagsKey], ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}
	return
}

func (i *Item) setLabel() (ok bool) {
	i.Label = i.Metadata[MetadataKeyForItemLabel(i.Type)]
	ok = i.Label != ""
	return
}

func (c *Catalogue) Tags() []string {
	tags := set.NewStringSet()
	for _, group := range c.Groups {
		for _, item := range group {
			for _, tag := range item.Tags() {
				tags.Insert(tag)
			}
		}
	}
	return tags.ToSlice()
}

func (c *Catalogue) HasMultipleGroups() bool {
	return len(c.Groups) > 1
}

// MakeCatalogue creates a Catalogue from the passed slice of items.
func MakeCatalogue(items []*Item) *Catalogue {
	if len(items) == 0 {
		return nil
	}

	catalogue := &Catalogue{Group(items)}
	for _, group := range catalogue.Groups {
		Sort(group)
	}
	return catalogue
}

// Group splits the passed slice into groups (new slices)
// of items having the same Type.
func Group(items []*Item) map[string][]*Item {
	if len(items) == 0 {
		return nil
	}

	groups := map[string][]*Item{}
	for _, item := range items {
		typeLabel := GroupTypeLabel(item.Type)
		groups[typeLabel] = append(groups[typeLabel], item)
	}

	return groups
}

// Sort sorts the passed slice of items as defined by item type.
// The function assumes that all items are of the same type.
// <#hardcoded#>
func Sort(items []*Item) {
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

func MetadataKeyForItemLabel(typeKey string) string {
	switch typeKey {
	case "books":
		return "title"
	case "games":
		return "title"
	default:
		return MDLabelKey
	}
}

// TypeLabel returns the label to be used in rendering of an item's type.
// <#hardcoded#>
func TypeLabel(typeKey string) string {
	switch typeKey {
	case "books":
		return "Book"
	case "games":
		return "Game"
	case "equipment":
		return "Equipment part"
	default:
		return "Inventory item"
	}
}

// TypeGroupLabel returns the label to be used in rendering
// of type name for a group of items.
// <#hardcoded#>
func GroupTypeLabel(typeKey string) string {
	switch typeKey {
	case "books":
		return "Books"
	case "games":
		return "Games"
	case "equipment":
		return "Equipment items"
	default:
		return typeKey
	}
}

// IsValidItemKeyWordForType checks if the key is a valid single
// item keyword in the database (e.g. "book" for "books" type).
// Keyword "item" is always valid for every type.
// <#hardcoded#>
func IsValidItemKeyWordForType(key string, typeKey string) bool {
	if key == "item" {
		return true
	}
	switch typeKey {
	case "books":
		return key == "book"
	case "games":
		return key == "game"
	}
	return false
}
