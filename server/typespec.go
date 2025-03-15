package server

// Sort sorts the passed slice of items as defined by item type.
// The function assumes that all items are of the same type.
func Sort(items []*Item) {
	if len(items) == 0 {
		return
	}

	switch items[0].Type {
	case "books":
		By(sortHintOrTitle).Sort(items)
	case "games":
		By(groupedUnderSeries).Sort(items)
	case "equipment":
		By(label).Sort(items)
	default:
		By(label).Sort(items)
	}
}

// ItemLabelKey returns which metadata key to use to lookup value
// used in rendering of an item's label.
func ItemLabelKey(typeKey string) string {
	switch typeKey {
	case "books":
		return "title"
	case "games":
		return "title"
	default:
		return MKEY_LABEL
	}
}

// TypeLabel returns the label to be used in rendering of an item's type name.
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
