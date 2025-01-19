package server

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

var (
	ErrNilReader     = errors.New("io.Reader is nil")
	ErrInvalidFormat = errors.New("invalid XML byte stream")
)

const (
	XMLNODE_KOIDATABASE  = "koidatabase"
	XMLATTR_CREATED      = "created"
	XMLATTR_LASTMODIFIED = "lastModified"
	XMLNODE_DATA         = "data"
	XMLNODE_METADATA     = "metadata"
	XMLNODE_KOITYPES     = "koitypes"
	XMLATTR_ENABLED      = "enabled"
	XMLNODE_COLLECTIONS  = "collections"
	XMLNODE_COLLECTION   = "collection"
	XMLATTR_HIDDEN       = "hidden"
)

func DecodeDatabase(r io.Reader) error {
	if r == nil {
		return ErrNilReader
	}

	decoder := xml.NewDecoder(r)

	// <koidatabase ... >
	currentNode, err := expectStart(decoder, XMLNODE_KOIDATABASE)
	if err != nil {
		return fmt.Errorf("failed to detect <%s>: %w", XMLNODE_KOIDATABASE, err)
	}
	created, _ := findAttribute(currentNode, XMLATTR_CREATED)
	lastModified, _ := findAttribute(currentNode, XMLATTR_LASTMODIFIED)
	if ts, err := time.Parse(time.DateOnly, created); err != nil {
		return fmt.Errorf("failed to detect or decode attribute <%s %s>: %w", XMLNODE_KOIDATABASE, XMLATTR_CREATED, err)
	} else {
		_database.created = ts
	}
	if ts, err := time.Parse(time.DateOnly, lastModified); err != nil {
		return fmt.Errorf("failed to detect or decode attribute <%s %s>: %w", XMLNODE_KOIDATABASE, XMLATTR_LASTMODIFIED, err)
	} else {
		_database.lastModified = ts
	}
	trace(_decoder, "created=%s, lastModified=%s", created, lastModified)

	// <koitypes ...>
	currentNode, err = expectStart(decoder, XMLNODE_KOITYPES)
	if err != nil {
		return fmt.Errorf("failed to detect <%s>: %w", XMLNODE_KOITYPES, err)
	}
	enabledTypesStr, found := findAttribute(currentNode, XMLATTR_ENABLED)
	if !found {
		return fmt.Errorf("failed to detect attribute <%s %s>", XMLNODE_KOITYPES, XMLATTR_ENABLED)
	}
	if enabledTypes := splitJoinedWords(enabledTypesStr); enabledTypes == nil {
		return fmt.Errorf("failed to decode attribute <%s %s>: invalid typelist format", XMLNODE_KOITYPES, XMLATTR_ENABLED)
	} else {
		for _, t := range enabledTypes {
			_database.enabledTypes.Insert(t)
		}
		trace(_decoder, "enabled types: %s", strings.Join(enabledTypes, ", "))
	}

	// <koitypes> <metadata ... /> 0..N </koitypes>
	for {
		if currentNode, err = nextOrEnd(decoder, XMLNODE_METADATA, XMLNODE_KOITYPES); err != nil {
			return fmt.Errorf("failed to detect <%s> or </%s>: %w", XMLNODE_METADATA, XMLNODE_KOITYPES, err)
		}
		if currentNode != nil {
			// <metadata ...>
			var metadata = &struct {
				Key          string `xml:"key,attr"`
				DefaultValue string `xml:"default,attr"`
			}{}
			if err = decoder.DecodeElement(metadata, currentNode); err != nil {
				return fmt.Errorf("failed to decode <%s>: %w", XMLNODE_METADATA, err)
			}
			if !isValidDefaultMetadataValueKeyRE(metadata.Key) {
				return fmt.Errorf("failed to decode <%s>: invalid attribute format", XMLNODE_METADATA)
			}
			_database.defaults[metadata.Key] = metadata.DefaultValue
			trace(_decoder, "predefined default [%s]->%s", metadata.Key, metadata.DefaultValue)
		} else {
			// </koitypes>
			trace(_decoder, "XML node <%s> decoding completed", XMLNODE_KOITYPES)
			break
		}
	}

	// <collections ...>
	currentNode, err = expectStart(decoder, XMLNODE_COLLECTIONS)
	if err != nil {
		return fmt.Errorf("failed to detect <%s>: %w", XMLNODE_COLLECTIONS, err)
	}
	hiddenCollectionsStr, found := findAttribute(currentNode, XMLATTR_HIDDEN)
	if found {
		if hiddenCollections := splitJoinedWords(hiddenCollectionsStr); hiddenCollections == nil {
			return fmt.Errorf("failed to decode attribute <%s %s>: invalid keylist format", XMLNODE_COLLECTIONS, XMLATTR_HIDDEN)
		} else {
			for _, c := range hiddenCollections {
				_database.hiddenCollections.Insert(c)
			}
			trace(_decoder, "hidden collections: %s", strings.Join(hiddenCollections, ", "))
		}
	}

	// <collections> <collection ... /> 0..N </collections>
	for {
		if currentNode, err = nextOrEnd(decoder, XMLNODE_COLLECTION, XMLNODE_COLLECTIONS); err != nil {
			return fmt.Errorf("failed to detect <%s> or </%s>: %w", XMLNODE_COLLECTION, XMLNODE_COLLECTIONS, err)
		}
		if currentNode != nil {
			// <collection ...>
			var collection = &struct {
				Key  string `xml:"key,attr"`
				Name string `xml:"name,attr"`
			}{}
			if err = decoder.DecodeElement(collection, currentNode); err != nil {
				return fmt.Errorf("failed to decode <%s>: %w", XMLNODE_COLLECTION, err)
			}
			if !isValidCollectionKey(collection.Key) {
				return fmt.Errorf("failed to decode <%s>: invalid attribute format", XMLNODE_COLLECTION)
			}
			_database.declaredCollections[collection.Key] = collection.Name
			trace(_decoder, "declared collection [%s]->%q", collection.Key, collection.Name)
		} else {
			// </collections>
			trace(_decoder, "XML node <%s> decoding completed", XMLNODE_COLLECTIONS)
			break
		}
	}

	// <data>
	if _, err = expectStart(decoder, XMLNODE_DATA); err != nil {
		return fmt.Errorf("failed to detect <%s>: %w", XMLNODE_DATA, err)
	}
	// <data> <TYPE>...</TYPE> 0..N </data>
	for {
		if currentNode, err = anyStartOrEnd(decoder, XMLNODE_DATA); err != nil {
			return fmt.Errorf("failed to detect <TYPE> or </%s>: %w", XMLNODE_DATA, err)
		}
		if currentNode != nil {
			// <TYPE>...</TYPE>
			typeKey := currentNode.Name.Local
			if !isValidWord(typeKey) {
				trace(_warning, "skipping XML node <%s> entirely: invalid typename format", typeKey)
				decoder.Skip()
				continue
			}
			if !_database.enabledTypes.Contains(typeKey) {
				trace(_decoder, "skipping XML node <%s> entirely: type is not enabled", typeKey)
				decoder.Skip()
				continue
			}
			trace(_decoder, "proceeding to decode XML node <%s> and all items defined for this type", typeKey)
			itemCnt := 0
			// <TYPE> <ITEM>...</ITEM> 0..N </TYPE>
			for {
				if currentNode, err = anyStartOrEnd(decoder, typeKey); err != nil {
					return fmt.Errorf("failed to detect <ITEM> or </%s>: %w", typeKey, err)
				}
				if currentNode != nil {
					// <ITEM>
					itemKey := currentNode.Name.Local
					if !IsValidItemKeyWordForType(itemKey, typeKey) {
						trace(_warning, "skipping XML node <%s> entirely: unknown keyword for items of type %q", itemKey, typeKey)
						decoder.Skip()
						continue
					}
					itemMetadata := make(map[string]string)
					for _, attr := range currentNode.Attr {
						if isValidMetadataKey(attr.Name.Local) {
							itemMetadata[attr.Name.Local] = attr.Value
						} else {
							trace(_warning, "skipping attribute <%s %s>: invalid metadata key format", itemKey, attr.Name.Local)
						}
					}
					if item := _database.AddItem(typeKey, itemMetadata); item == nil {
						trace(_warning, "failed to add item of type %q to the database, check item metadata", typeKey)
						decoder.Skip()
						continue
					}
					itemCnt++
					decoder.Skip()
				} else {
					// </TYPE>
					trace(_decoder, "decoded %d items of type %q", itemCnt, typeKey)
					break
				}
			}
		} else {
			// </data>
			trace(_decoder, "XML node <%s> decoding completed", XMLNODE_DATA)
			break
		}
	}

	// </koidatabase>
	if _, err = expectEnd(decoder, XMLNODE_KOIDATABASE); err != nil {
		return fmt.Errorf("failed to detect </%s>: %w", XMLNODE_KOIDATABASE, err)
	}

	return nil
}

func expectStart(decoder *xml.Decoder, name string) (*xml.StartElement, error) {
	tok, err := nextToken(decoder)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	switch tok := tok.(type) {
	case xml.StartElement:
		if tok.Name.Local == name {
			return &tok, nil
		}
	}
	return nil, ErrInvalidFormat
}

func expectEnd(decoder *xml.Decoder, name string) (*xml.EndElement, error) {
	tok, err := nextToken(decoder)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	switch tok := tok.(type) {
	case xml.EndElement:
		if tok.Name.Local == name {
			return &tok, nil
		}
	}
	return nil, ErrInvalidFormat
}

func nextOrEnd(decoder *xml.Decoder, next string, end string) (*xml.StartElement, error) {
	tok, err := nextToken(decoder)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	switch tok := tok.(type) {
	case xml.StartElement:
		if tok.Name.Local == next {
			return &tok, nil
		}
	case xml.EndElement:
		if tok.Name.Local == end {
			return nil, nil
		}
	}
	return nil, ErrInvalidFormat
}

func anyStartOrEnd(decoder *xml.Decoder, end string) (*xml.StartElement, error) {
	tok, err := nextToken(decoder)
	if err != nil {
		return nil, ErrInvalidFormat
	}
	switch tok := tok.(type) {
	case xml.StartElement:
		return &tok, nil
	case xml.EndElement:
		if tok.Name.Local == end {
			return nil, nil
		}
	}
	return nil, ErrInvalidFormat
}

func findAttribute(elem *xml.StartElement, attrName string) (val string, ok bool) {
	for _, attr := range elem.Attr {
		if attr.Name.Local == attrName {
			ok = true
			val = attr.Value
			return
		}
	}
	return
}

func nextToken(decoder *xml.Decoder) (tok xml.Token, err error) {
	for {
		tok, err = decoder.Token()
		if err != nil {
			return
		}
		switch t := tok.(type) {
		case xml.CharData:
			if len(strings.TrimSpace(string(t))) > 0 {
				return
			}
		default:
			return
		}
	}
}
