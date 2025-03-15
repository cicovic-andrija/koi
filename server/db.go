package server

import (
	"strings"
	"time"

	"src.acicovic.me/koipond/set"
)

// Special metadata keys.
const (
	MKEY_LABEL        string = "label"
	MKEY_COLLECTIONS  string = "collections"
	MKEY_TAGS         string = "tags"
	MKEY_SORTING_HINT string = "sortBy"
)

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
	metadata = nil

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

	// index collection (invalid and hidden collections are cleaned out)
	if collections := item.Metadata[MKEY_COLLECTIONS]; collections != "" {
		validCollections := []string{}
		for _, collectionKey := range strings.Split(collections, ",") {
			collectionKey = strings.TrimSpace(collectionKey)
			if isValidCollectionKey(collectionKey) {
				_, declared := db.declaredCollections[collectionKey]
				if declared && !db.hiddenCollections.Contains(collectionKey) {
					validCollections = append(validCollections, collectionKey)
				}
			}
		}
		for _, c := range validCollections {
			db.collectioned[c] = append(db.collectioned[c], item)
		}
		item.Metadata[MKEY_COLLECTIONS] = strings.Join(validCollections, ",")
	}

	// index tags (invalid tags are cleaned out)
	if tags := item.Metadata[MKEY_TAGS]; tags != "" {
		validTags := []string{}
		for _, tag := range strings.Split(tags, ",") {
			tag = strings.TrimSpace(tag)
			if isValidTag(tag) {
				validTags = append(validTags, tag)
			}
		}
		for _, tag := range validTags {
			db.tagged[tag] = append(db.tagged[tag], item)
		}
		item.Metadata[MKEY_TAGS] = strings.Join(validTags, ",")
	}

	db.items = append(db.items, item)
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

func (db *Database) catalogueOfEverything() *Catalogue {
	return makeCatalogue(db.items)
}

func (db *Database) catalogueForCollection(key string) *Catalogue {
	return makeCatalogue(db.collectioned[key])
}

func (db *Database) catalogueOfTaggedItems(tag string) *Catalogue {
	return makeCatalogue(db.tagged[tag])
}

func makeCatalogue(items []*Item) *Catalogue {
	if len(items) == 0 {
		return nil
	}

	catalogue := &Catalogue{
		groups: group(items),
	}
	for _, group := range catalogue.groups {
		Sort(group)
	}

	return catalogue
}

func group(items []*Item) map[string][]*Item {
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
