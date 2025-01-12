package server

import "src.acicovic.me/koipond/set"

var db = &Database{
	items:        []*Item{},
	defaults:     map[string]string{},
	enabledTypes: set.NewStringSet(),
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
