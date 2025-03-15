package server

import (
	"sort"
	"strings"

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
// e.g. the Sort function.
//
// Item implements CommonDataProperties.
type Item struct {
	CommonBaseObject

	ID       int
	Type     string
	Label    string
	Metadata map[string]string
}

// Catalogue is a collection of items grouped by type.
//
// Catalogue implements CommonDataProperties.
type Catalogue struct {
	CommonBaseObject

	groups   map[string][]*Item
	hideTags bool
}

// TagMap implements CommonDataProperties.
type TagMap struct {
	CommonBaseObject

	ref map[string]int
}

// CollectionMap implements CommonDataProperties.
type CollectionMap struct {
	CommonBaseObject

	ref map[string]string
}

func (i *Item) setLabel() (ok bool) {
	i.Label = i.Metadata[ItemLabelKey(i.Type)]
	ok = i.Label != ""
	return
}

// Ref implements CommonDataProperties.
func (i *Item) Ref() any {
	return i
}

// Properties implements CommonDataProperties.
func (i *Item) Properties() map[string]string {
	return i.Metadata
}

// Tags returns a slice of item's tags.
func (i *Item) Tags() (tags []string) {
	if i.Metadata[MKEY_TAGS] != "" {
		tags = strings.Split(i.Metadata[MKEY_TAGS], ",")
		sort.Strings(tags)
	}
	return
}

// Tags returns a slice of tags found in the catalogue.
func (c *Catalogue) Tags() (tags []string) {
	tagset := set.NewStringSet()
	for _, group := range c.groups {
		for _, item := range group {
			for _, tag := range item.Tags() {
				tagset.Insert(tag)
			}
		}
	}
	tags = tagset.ToSlice()
	sort.Strings(tags)
	return
}

// HideTags implements CommonDataProperties.
func (c *Catalogue) HideTags() bool {
	return c.hideTags
}

func (c *Catalogue) withHiddenTags() *Catalogue {
	c.hideTags = true
	return c
}

// Groups implements CommonDataProperties.
//
// This is a bit of an expensive function, since Go doesn't allow conversion between []*Ptr to []interface{},
// or in this case []*Item to []CommonDataProperties, so conversion needs to be done for each element in every group.
func (c *Catalogue) Groups() map[string][]CommonDataProperties {
	groups := make(map[string][]CommonDataProperties)
	for k, v := range c.groups {
		slice := make([]CommonDataProperties, len(v))
		for j, i := range v {
			slice[j] = i
		}
		groups[k] = slice
	}
	return groups
}

// MultiGroup implements CommonDataProperties.
func (c *Catalogue) MultiGroup() bool {
	return len(c.groups) > 1
}

// HasMultipleGroups reports whether a Catalogue has more than one group.
func (c *Catalogue) HasMultipleGroups() bool {
	return len(c.groups) > 1
}

// NewTagMap return a pointer to a new TagMap object.
func NewTagMap(ref map[string]int) *TagMap {
	return &TagMap{ref: ref}
}

// Ref implements CommonDataProperties.
func (t *TagMap) Ref() any {
	return t.ref
}

// HideTags implements CommonDataProperties.
func (*TagMap) HideTags() bool {
	return true
}

// NewCollectionMap return a pointer to a new CollectionMap object.
func NewCollectionMap(ref map[string]string) *CollectionMap {
	return &CollectionMap{ref: ref}
}

// Ref implements CommonDataProperties.
func (c *CollectionMap) Ref() any {
	return c.ref
}

// HideTags implements CommonDataProperties.
func (*CollectionMap) HideTags() bool {
	return true
}
