package fansischmancy

import (
	"bytes"
	"io"
)

// NonAnsiDetectionWriter wraps an io.Writer and detects non-4-bit ANSI color codes.
type NonAnsiDetectionWriter struct {
	w io.Writer
}

// NewFANSISchmancy creates a new writer that detects non-4-bit ANSI color codes.
func NewWriter(w io.Writer) *NonAnsiDetectionWriter {
	return &NonAnsiDetectionWriter{w: w}
}

// Write implements the io.Writer interface.
func (w *NonAnsiDetectionWriter) Write(p []byte) (n int, err error) {
	var output []byte
	parts := bytes.Split(p, []byte("\x1b["))

	// First part is before any escape sequence
	output = append(output, parts[0]...)

	for i, part := range parts[1:] {
		if len(part) == 0 {
			// Empty part means we had a bare escape sequence
			output = append(output, '\x1b', '[')
			continue
		}

		// Look for the end of the sequence (m)
		if idx := bytes.IndexByte(part, 'm'); idx >= 0 {
			sequence := part[:idx]
			remainder := part[idx+1:]

			// Always write the original sequence
			output = append(output, '\x1b', '[')
			output = append(output, part[:idx+1]...)

			// Find any newline in the remainder
			if nlIdx := bytes.IndexByte(remainder, '\n'); nlIdx >= 0 {
				// Write text up to newline
				if !isSimpleColorCode(sequence) && (bytes.HasPrefix(sequence, []byte("38;2")) ||
					bytes.HasPrefix(sequence, []byte("48;2")) ||
					bytes.HasPrefix(sequence, []byte("38;5")) ||
					bytes.HasPrefix(sequence, []byte("48;5"))) {
					output = append(output, '\x1b', '[', '7', ';', '9', 'm') // reverse video and strikethrough
					output = append(output, remainder[:nlIdx]...)
					output = append(output, '\x1b', '[', '2', '7', ';', '2', '9', 'm') // reset reverse video and strikethrough
				} else {
					output = append(output, remainder[:nlIdx]...)
				}
				// Write newline and rest of text
				output = append(output, remainder[nlIdx:]...)
			} else {
				// No newline, write remainder with effects if needed
				if !isSimpleColorCode(sequence) && (bytes.HasPrefix(sequence, []byte("38;2")) ||
					bytes.HasPrefix(sequence, []byte("48;2")) ||
					bytes.HasPrefix(sequence, []byte("38;5")) ||
					bytes.HasPrefix(sequence, []byte("48;5"))) {
					output = append(output, '\x1b', '[', '7', ';', '9', 'm') // reverse video and strikethrough
					output = append(output, remainder...)
					output = append(output, '\x1b', '[', '2', '7', ';', '2', '9', 'm') // reset reverse video and strikethrough
				} else {
					output = append(output, remainder...)
				}
			}
		} else {
			// No 'm' found, just write the part as-is with the escape sequence
			output = append(output, '\x1b', '[')
			if i == len(parts[1:])-1 {
				// Last part and no 'm', treat as incomplete
				output = append(output, part...)
			} else {
				// Not the last part, check for 24-bit or 256 color
				if bytes.HasPrefix(part, []byte("38;2")) ||
					bytes.HasPrefix(part, []byte("48;2")) ||
					bytes.HasPrefix(part, []byte("38;5")) ||
					bytes.HasPrefix(part, []byte("48;5")) {
					output = append(output, '\x1b', '[', '7', ';', '9', 'm') // reverse video and strikethrough
					output = append(output, part...)
					output = append(output, '\x1b', '[', '2', '7', ';', '2', '9', 'm') // reset reverse video and strikethrough
				} else {
					output = append(output, part...)
				}
			}
		}
	}

	_, err = w.w.Write(output)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func isSimpleColorCode(sequence []byte) bool {
	// Check if sequence is just a number between 30-37, 40-47, or 0
	if len(sequence) == 0 {
		return false
	}

	// Handle reset code
	if bytes.Equal(sequence, []byte("0")) {
		return true
	}

	// Split on semicolon for bright/bold variants
	parts := bytes.Split(sequence, []byte(";"))
	if len(parts) > 2 {
		return false
	}

	// Check first part is a valid color code
	if len(parts[0]) != 2 {
		return false
	}
	first := parts[0][0]
	second := parts[0][1]
	isColor := (first == '3' || first == '4') && second >= '0' && second <= '7'

	// If there's a second part, it should be "1" for bright/bold
	if len(parts) == 2 {
		return isColor && bytes.Equal(parts[1], []byte("1"))
	}

	return isColor
}
