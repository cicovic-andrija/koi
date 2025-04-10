# koipond: personal inventory server

`v1.x`, `self-hosted`, `minimal`, `do-it-yourself`, `web-server`, `docker-containerized`

# 1. About

`koipond` is a software system for personal inventory management. The north star of the project is
to provide a simple and minimal web interface that can be used to browse and update small-scale
collections of items (thousands of items), with raw data persisted in textual, human-readable,
standard format. It is designed to be easily extensible, customizable and self-hosted, for users
that are comfortable with technology and prefer minimal systems.

In the current version `1.x`, the server offers only browsing (read-only) functionality, meaning
that the XML file used for persistent storage must be populated/updated through a text editor program.
For the one person currently working on this project, that was a high enough bar for version `1.x` of
the system (demo link below). This decision also aligns with the idea of a target audience for this
system. Version `2.x` will offer read-write capabilities (to be implemented).

### Version `1.x` Demo: [https://inventory.acicovic.me](https://inventory.acicovic.me)

# 2. Before Further Reading

For the rest of this document, notes related to source code pointers and implementation details
will be visually separated by using Markdown footnote formatting, for example:

> The text you are currently reading can be found in the `README.md` at the root of the source tree.

# 3. Concepts

- Acting entities in the software system are SERVER, CLIENT and USER.
- USER utilizes a CLIENT program (e.g. a web browser) to send requests and present responses from the SERVER.
- CLIENT and SERVER communicate over a standard network protocol (HTTPS) to exchange requests and responses.
  - In simpler words, this is a website.
- SERVER uses a textual (XML) file to persist data over time, called just the database in further text.
- Database is designed to be human-readable; this eliminates the need to implement import/export mechanisms.
- SERVER manages generic items that are stored in the database.
- SERVER is responsible that every item is uniquely identifiable.
- USER provides details for each item: a set of (meta)data key-value pairs, called just metadata in further text.
  - Example: for a book (item), USER provides metadata keys like `title`, `author`, `edition` and their corresponding values.
- Both keys and values in metadata are always interpreted as text (UTF-8 encoded strings) by the SERVER.
- The only required metadata key-value pair is that which determines how an item will be labeled (item's "name").
- SERVER ignores items for which a label cannot be determined from the provided metadata set.
- Item may belong to one or more collections, specified by the USER with a special metadata key that SERVER is able to detect.
  - Example: a book (item) with label `"The Silmarillion"` may belong to a collection `"Tolkien's Legendarium"`.
- Item may be tagged with one or more tags, specified by the USER with a special metadata key that SERVER is able to detect.
  - Example: a book (item) with label `"Rendezvous With Rama"` may have tags `novel` and `scifi`.
- Each item is of one and only one type, specified by the USER.
  - Example: a book (item) is of type `books`.
- Type is a classification used to logically partition items when there is a need to perform an operation on many items.
  - SERVER always groups items by their type in listings.
  - SERVER always sorts items of the same type in a type-specific way.
  - In other words, a type defines "behavior" of the SERVER when dealing with items of that type.
- Types are generic and have well-defined default behavior, but it is possible to customize types with specific behavior.
  - Example: for a generic `books` type the SERVER would look for a `label` metada key to determine a book's label,
    however the type can be customized so the SERVER looks for a `title` metadata key to determine a book's label.
  - Example: for a generic `books` type the SERVER would sort a list of books by their labels,
    however the type can be customized so the SERVER uses metadata value for key `sortBy` to sort a list of books.
- A collection may be composed of items of different types.
- Items of different types can be tagged with same tags.

# 4. Expected Format

This section describes the format of XML files that the system is able to parse. An example file can be found in
`examples/koidata.xml`. Placeholders that need to be replaced are denoted by all caps `LIKE-THIS`.

## 4.1. Root

```xml
<koidatabase created="YYYY-MM-DD" lastModified="YYYY-MM-DD">
```

The root XML node that is expected at the beginning of the file is the `<koidatabase>` node;
it should be the parent of all other nodes that the parser is scanning for. In version `1.x`,
the end user is responsible for the `created` and `lastModified` attributes.

> **Important note: The parser is not able to recognize and ignore XML comments (will be fixed).**

## 4.2. Enabled Types and Default Values

```xml
<koitypes enabled="TYPENAME-1,TYPENAME-2,TYPENAME3"> <!-- Child of <koidatabase> -->
  <metadata key="TYPENAME-1/KEY-A" default="TYPE-1-DEFAULT-VALUE-A" />
  <metadata key="TYPENAME-2/KEY-A" default="TYPE-2-DEFAULT-VALUE-A" />
  <metadata key="TYPENAME-2/KEY-B" default="TYPE-2-DEFAULT-VALUE-B" />
  <metadata key="TYPENAME-3/KEY-B" default="TYPE-3-DEFAULT-VALUE-B" />
</koitypes>
```

The immediate child of the root node should be the `<koitypes>` node. The `enabled` attribute
should specify which item types will be considered during item scanning later on. If there are
items of a type not listed here, the server will ignore them. This is useful for tracking
items in the database for tracking purposes only, as they will not be loaded into the system.

> **Important note: Type names are composed only of lowercase English characters a-z.**

The `<koitypes>` node can optionally have one or more `<metadata>` child nodes which should
specify default metadata values for items that will not explicitly declare values for given
keys. An attribute in format `key="TYPENAME-1/KEY-A"` designates all items of type `TYPENAME-1`
that have not explicitly declared a value for key `KEY-A`. They will be assigned the value
given in the `default` attribute for key `KEY-A`, implicitly.

## 4.3. Collections

```xml
<collections> <!-- Child of <koidatabase> -->
  <collection key="COLLECTION-KEY-1" name="COLLECTION-NAME-1" />
  <collection key="COLLECTION-KEY-2" name="COLLECTION-NAME-2" />
</collections>
```

## 4.4. Items Grouped by Types

```xml
<data> <!-- Child of <koidatabase> -->
  <TYPENAME-1>
    ... <!-- Items of type TYPENAME-1 -->
  </TYPENAME-1>
  <TYPENAME-2>
    ... <!-- Items of type TYPENAME-2 -->
  </TYPENAME-2>
</data>
```

## 4.5. Items

```xml
<TYPENAME-1> <!-- Child of <data> -->
  <item label="ITEM-LABEL" KEY-1="VAL-1" KEY-2="VAL-2" ... />
  <!-- OR -->
  <TYPE-1-SPECIFIC-KEYWORD label="ITEM-LABEL" KEY-1="VAL-1" KEY-2="VAL-2" ... />
  <!-- Mixing of this two formats is allowed -->
</TYPENAME-2>
```

### 4.5.1. Additional Notes About Items and Types
> **Important note: Label can also be determined from other key.**

> **Important note: Collections.**

> **Important note: Tags.**

## 4.6. Virtual End of File

```xml
</koidatabase> <!-- EOF -->
```

Parser does not attempt to look for tokens past this point in the file - it considers the closing
token of `</koidatabase>` as EOF, even if the file has content after the token.

# 5. How To: Programming, Building, Testing

### Build

```bash
$ go build -o koipond main.go
```

### Run locally (development mode)

TODO: Mention `store/`

```bash
$ KOIPOND_MODE=dev go run main.go
```

# 6. How To: Deployment

> TCP port in dev mode is hard-coded to 8072.

### Run in production

```bash
$ KOIPOND_MODE=prod-local-listener KOIPOND_PORT=52000 ./koipond
```

> Requires a service manager to handle crashes and log redirection. See `systemd.service` for an example.

> For encrypted traffic, configure a reverse HTTPS proxy, e.g. `nginx`.

> For authentication, configure a stanalone authentication service.

### Run in production (Docker)

TODO: Mention store/

```bash
$ docker run \
    --name inventory-server \
    --publish 127.0.0.1:8072:8072 \
    --env 'KOIPOND_PORT=8072' \
    --volume $HOME/store:/srv/store \
    --restart on-failure:10 \
    --detach \
    acicovic/koipond:latest
```

### Kill

```bash
$ pkill -SIGINT koipond
```

> Or, CTRL-C if running in foreground.

or

```bash
$ docker kill -s SIGINT koipond-server
```

### Build Docker image

```bash
$ docker build -t acicovic/koipond:v1.x .
$ docker tag acicovic/koipond:v1.x acicovic/koipond:latest
$ docker push acicovic/koipond:v1.x
$ docker push acicovic/koipond:latest
```

> Substitute image name.
