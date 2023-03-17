package image

import (
	"errors"
	"image"
	"image/jpeg"
	"io"
)

// jfifChunk is a sequence: an SOI chunk, an APP0/JFIF chunk and finally the
// 0xff that starts the third chunk.
var jfifChunk = []byte{
	0xff, 0xd8, // SOI  marker.
	0xff, 0xe0, // APP0 marker.
	0x00, 0x10, // Length: 16 byte payload (including these two bytes).
	0x4a, 0x46, 0x49, 0x46, 0x00, // "JFIF\x00".
	0x01, 0x01, // Version 1.01.
	0x00,       // No density units.
	0x00, 0x01, // Horizontal pixel density.
	0x00, 0x01, // Vertical   pixel density.
	0x00, // Thumbnail width.
	0x00, // Thumbnail height.
	0xff, // Start of the third chunk's marker.
}

func jfifEncode(w io.Writer, m image.Image, o *jpeg.Options) error {
	return jpeg.Encode(&jfifWriter{w: w}, m, o)
}

// jfifWriter wraps an io.Writer to convert the data written to it from a plain
// JPEG to a JFIF-enhanced JPEG. It implicitly buffers the first three bytes
// written to it. The fourth byte will tell whether the original JPEG already
// has the APP0 chunk that JFIF requires.
type jfifWriter struct {
	// w is the wrapped io.Writer.
	w io.Writer
	// n ranges between 0 and 4 inclusive. It is the number of bytes written to
	// this (which also implements io.Writer), saturating at 4. The first three
	// bytes are expected to be {0xff, 0xd8, 0xff}. The fourth byte indicates
	// whether the second JPEG chunk is an APP0 chunk or something else.
	n int
}

func (jw *jfifWriter) Write(p []byte) (int, error) {
	nSkipped := 0

	for jw.n < 3 {
		if len(p) == 0 {
			return nSkipped, nil
		} else if p[0] != jfifChunk[jw.n] {
			return nSkipped, errors.New("jfifWriter: input was not a JPEG")
		}
		nSkipped++
		jw.n++
		p = p[1:]
	}

	if jw.n == 3 {
		if len(p) == 0 {
			return nSkipped, nil
		}
		chunk := jfifChunk
		if p[0] == 0xe0 {
			// The input JPEG already has an APP0 marker. Just write SOI (2
			// bytes) and an 0xff: the three bytes we've previously skipped.
			chunk = chunk[:3]
		}
		if _, err := jw.w.Write(chunk); err != nil {
			return nSkipped, err
		}
		jw.n = 4
	}

	n, err := jw.w.Write(p)
	return n + nSkipped, err
}
