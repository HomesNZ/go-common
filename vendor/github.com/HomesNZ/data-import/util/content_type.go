package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// ContentType reads the first 512 bytes of a file to determine it's mime content type. If the first 512 bytes cannot be
// read or the content type cannot be determined, it returns "application/octet-stream"
func ContentType(r io.ReadSeeker) string {
	// we should only be peeking here rather than reading as we miss this data later.
	buff := make([]byte, 512)
	_, err := r.Read(buff)
	// Make sure we rewind the reader back to the start.
	r.Seek(0, 0)
	if err != nil {
		fmt.Println("Error determining content type: ", err)
		return "application/octet-stream"
	}

	buff = trapBOM(buff)
	contentType := http.DetectContentType(buff)

	// BoxAndDice are sending xml files with no xml prolog, so the DetectContentType function returns text/plain, which we then
	// refuse to interpret. This interrogates the XML to see if it starts with <propertyList, in which case it will return as
	// text/xml. The check (<propertyList) may need to be updated if there are other XML files in the same situation, but at the moment
	// this is confined to Reaxml (and BoxAndDice).
	if strings.Contains(contentType, "text/plain") {
		firstThirteenCharacters := string(buff[0:13])
		if strings.Contains(firstThirteenCharacters, "<propertyList") {
			contentType = strings.Replace(contentType, "plain", "xml", -1)
		}
	}
	return contentType
}

// TrapBOM removes the Byte Order Mark (https://en.wikipedia.org/wiki/Byte_order_mark). It indicates that the file is a UTF-8 file.
// FirstNational example reaxml files use this, which fails the http.DetectContentType test, so we need to remove it.
func trapBOM(fileBytes []byte) []byte {
	trimmedBytes := bytes.Trim(fileBytes, "\xef\xbb\xbf")
	return trimmedBytes
}
