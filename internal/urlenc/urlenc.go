package urlenc

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io"
	"sort"

	"oss.terrastruct.com/d2/d2graph"
)

var compressionDict = "->" +
	"<-" +
	"--" +
	"<->"

func init() {
	var common []string
	for k := range d2graph.StyleKeywords {
		common = append(common, k)
	}
	for k := range d2graph.ReservedKeywords {
		common = append(common, k)
	}
	for k := range d2graph.ReservedKeywordHolders {
		common = append(common, k)
	}
	sort.Strings(common)
	for _, k := range common {
		compressionDict += k
	}
}

// Decode decodes a compressed base64 D2 string.
func Decode(encoded string) (_ string, err error) {
	b64Decoded, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	zr := flate.NewReaderDict(bytes.NewReader(b64Decoded), []byte(compressionDict))
	var b bytes.Buffer
	if _, err := io.Copy(&b, zr); err != nil {
		return "", err
	}
	if err := zr.Close(); err != nil {
		return "", nil
	}
	return b.String(), nil
}
