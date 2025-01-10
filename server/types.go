package server

type Item struct {
	t *ItemType

	ID       int               `json:"id"`
	Label    string            `json:"label"`
	Metadata map[string]string `json:"metadata"`
}

type Catalogue struct {
	Groups map[string][]*Item `json:"groups"`
}

type ItemType struct {
	key string
}

func (i *Item) Type() string {
	return i.t.key
}
