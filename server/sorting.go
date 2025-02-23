package server

import "sort"

// By is the type of a "less" function that defines the ordering of items.
type By func(*Item, *Item) bool

// Type sorter implements sort.Interface.
type sorter struct {
	items []*Item
	by    func(*Item, *Item) bool
}

// Sort is a method on the function type By that sorts the slice as defined by the function.
func (by By) Sort(items []*Item) {
	sort.Sort(&sorter{items, by})
}

// Len is part of sort.Interface.
func (s *sorter) Len() int {
	return len(s.items)
}

// Swap is part of sort.Interface.
func (s *sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

// Less is part of sort.Interface.
func (s *sorter) Less(i, j int) bool {
	return s.by(s.items[i], s.items[j])
}

func label(i *Item, j *Item) bool {
	return i.Label < j.Label
}

func sortHintOrTitle(i *Item, j *Item) bool {
	const MKEY_TITLE = "title"
	a := i.Metadata[MKEY_SORTING_HINT]
	if a == "" {
		a = i.Metadata[MKEY_TITLE]
	}
	b := j.Metadata[MKEY_SORTING_HINT]
	if b == "" {
		b = j.Metadata[MKEY_TITLE]
	}
	return a < b
}

func groupedUnderSeries(i *Item, j *Item) bool {
	const (
		MKEY_TITLE  = "title"
		MKEY_SERIES = "series"
	)
	a := i.Metadata[MKEY_SERIES]
	if a == "" {
		a = i.Metadata[MKEY_TITLE]
	}
	b := j.Metadata[MKEY_SERIES]
	if b == "" {
		b = j.Metadata[MKEY_TITLE]
	}
	return a < b
}
