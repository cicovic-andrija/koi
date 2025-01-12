package server

import (
	"errors"
	"io"
)

var (
	ErrNilReader = errors.New("io.Reader is nil")
)

const (
	XMLNODE_TIMESTAMP = "timestamp"
	XMLNODE_DATA      = "data"
	XMLNODE_METADATA  = "metadata"
	XMLNODE_KOITYPES  = "koitypes"
)

type XMLNodeTimestamp struct {
	Created      string `xml:"created,attr"`
	LastModified string `xml:"lastModified,attr"`
}

type XMLNodeKoiTypes struct {
	Enabled string `xml:"enabled,attr"`
}

type XMLNodeMetadata struct {
	Key          string `xml:"key,attr"`
	DefaultValue string `xml:"default,attr"`
}

func DecodeDatabase(r io.Reader) error {
	if r == nil {
		return ErrNilReader
	}

	return nil
}
