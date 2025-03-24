package server

import (
	"net/http"
	"os"
	"strconv"
)

func multiHandler() http.Handler {
	// multiplexing handler
	mux := http.NewServeMux()

	register := func(p string, h func(http.ResponseWriter, *http.Request)) {
		mux.HandleFunc(p, h)
		trace(_https, "handler registered for pattern %s", p)
	}

	register("GET /collections", renderCollections)
	register("GET /collections/{collection}", renderCollection)
	// /collections/{$} -> 404
	register("GET /tags", renderTags)
	register("GET /tags/{tag}", renderTag)
	// /tags/{$} -> 404
	register("GET /items", renderItems)
	register("GET /items/{id}", renderItem)
	// /items/{$} 404
	register("GET /", defaultHandler)

	return mux
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "/collections", http.StatusMovedPermanently)
	case "/favicon.ico":
		var (
			icon *os.File
			fi   os.FileInfo
		)
		icon, err := os.Open("data/favicon.ico")
		if err == nil {
			fi, err = icon.Stat()
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		http.ServeContent(w, r, "favicon.ico", fi.ModTime(), icon)
	default:
		renderNotFound(w, "Page not found.")
	}
}

func renderCollections(w http.ResponseWriter, r *http.Request) {
	render(
		w,
		HTMLPage{
			Key:        "@collections",
			Supertitle: "All",
			Title:      "Collections",
			Data:       NewCollectionMap(_database.collections()),
		},
	)
}

func renderCollection(w http.ResponseWriter, r *http.Request) {
	collectionKey := r.PathValue("collection")
	catalogue := _database.catalogueForCollection(collectionKey)
	if catalogue == nil {
		renderNotFound(w, "Collection not found.")
		return
	}

	render(
		w,
		HTMLPage{
			Key:        "@catalogue",
			Supertitle: "Collection",
			Title:      _database.declaredCollections[collectionKey],
			Data:       catalogue,
		},
	)
}

func renderTags(w http.ResponseWriter, r *http.Request) {
	render(
		w,
		HTMLPage{
			Key:        "@tags",
			Supertitle: "All",
			Title:      "Tags",
			Data:       NewTagMap(_database.tags()),
		},
	)
}

func renderTag(w http.ResponseWriter, r *http.Request) {
	tag := r.PathValue("tag")
	catalogue := _database.catalogueOfTaggedItems(tag)
	if catalogue == nil {
		renderNotFound(w, "Tag not found.")
		return
	}

	render(
		w,
		HTMLPage{
			Key:        "@catalogue",
			Supertitle: "Items tagged with",
			Title:      tag,
			Data:       catalogue.withHiddenTags(),
		},
	)
}

func renderItems(w http.ResponseWriter, r *http.Request) {
	render(
		w,
		HTMLPage{
			Key:        "@catalogue",
			Supertitle: "All",
			Title:      "Items",
			Data:       _database.catalogueOfEverything().withHiddenTags(),
		},
	)
}

func renderItem(w http.ResponseWriter, r *http.Request) {
	itemID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || itemID < 0 || itemID > _database.lastID() {
		renderNotFound(w, "Item not found.")
		return
	}

	item := _database.singleItem(itemID)
	render(
		w,
		HTMLPage{
			Key:        "@" + item.Type + "/item",
			Supertitle: TypeLabel(item.Type),
			Title:      item.Label,
			Data:       item,
		},
	)
}

func renderNotFound(w http.ResponseWriter, message string) {
	render(
		w,
		HTMLPage{
			Key:          "@not-found",
			Supertitle:   "404",
			Title:        "Not Found",
			ErrorMessage: message,
			Data:         &CommonBaseObject{},
		},
	)
}
