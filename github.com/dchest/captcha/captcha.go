package captcha

import (
	"bytes"
	"io"
	"encoding/base64"
)

const (
	// Default number of digits in captcha solution.
	DefaultLen = 4
)

// Creator: chending
// NewId creates a new captcha with the standard length, saves it in the internal
// storage and returns its id and its digit.
func New() (string, string) {
	return NewLen(DefaultLen)
}


func NewLen(length int) (result, digit string) {
	digitBytes := RandomDigits(length)
	var content bytes.Buffer
	WriteImage(&content, digitBytes, StdWidth, StdHeight)
	result = base64.StdEncoding.EncodeToString(content.Bytes())
	for _, d := range digitBytes {
		digit += string('0' + d)
	}
	return
}


// WriteImage writes PNG-encoded image representation of the captcha with the
// given id. The image will have the given width and height.
func WriteImage(w io.Writer, digit []byte, width, height int) error {
	_, err := NewImage(digit, width, height).WriteTo(w)
	return err
}