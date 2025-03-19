package fansischmancy_test

import (
	"bytes"
	"testing"

	"github.com/williammartin/fansischmancy"
)

func TestNonAnsiDetectionWriter_SimpleWrite(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	input := []byte("Hello, World!\n")
	n, err := writer.Write(input)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if n != len(input) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(input), n)
	}

	if got := buf.String(); got != "Hello, World!\n" {
		t.Errorf("Expected 'Hello, World!\n', got %q", got)
	}
}

func TestNonAnsiDetectionWriter_Detect24BitColor(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// This is a 24-bit color code (RGB: 255, 100, 0)
	input := []byte("\x1b[38;2;255;100;0mColored Text\x1b[0m\n")
	_, err := writer.Write(input)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "\x1b[38;2;255;100;0m\x1b[9mColored Text\x1b[29m\x1b[0m\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}

func TestNonAnsiDetectionWriter_Allow4BitColor(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// These are 4-bit color codes that should be allowed
	input := []byte("\x1b[31mRed\x1b[0m \x1b[44mBlue BG\x1b[0m\n")
	n, err := writer.Write(input)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(input) {
		t.Errorf("Expected to write %d bytes, wrote %d", len(input), n)
	}

	// Should pass through without any detection message
	if got := buf.String(); got != string(input) {
		t.Errorf("Expected:\n%q\nGot:\n%q", input, got)
	}
}

func TestNonAnsiDetectionWriter_Multiple24BitSequences(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// Multiple 24-bit color sequences
	input := []byte("\x1b[38;2;255;100;0mOrange\x1b[0m \x1b[38;2;100;255;100mGreen\x1b[0m\n")
	_, err := writer.Write(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "\x1b[38;2;255;100;0m\x1b[9mOrange\x1b[29m\x1b[0m \x1b[38;2;100;255;100m\x1b[9mGreen\x1b[29m\x1b[0m\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}

func TestNonAnsiDetectionWriter_MultipleWrites(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// First write
	input1 := []byte("\x1b[38;2;255;100;0mOrange\n")
	_, err := writer.Write(input1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Second write
	input2 := []byte("\x1b[38;2;100;255;100mGreen\x1b[0m\n")
	_, err = writer.Write(input2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "\x1b[38;2;255;100;0m\x1b[9mOrange\x1b[29m\n\x1b[38;2;100;255;100m\x1b[9mGreen\x1b[29m\x1b[0m\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}

func TestNonAnsiDetectionWriter_EmptyWrite(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	_, err := writer.Write([]byte{})
	if err != nil {
		t.Errorf("Expected no error on empty write, got %v", err)
	}

	if got := buf.String(); got != "" {
		t.Errorf("Expected empty string, got %q", got)
	}
}

func TestNonAnsiDetectionWriter_IncompleteSequence(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// Incomplete ANSI sequence without 'm'
	input := []byte("\x1b[38;2;255;100;0")
	_, err := writer.Write(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := buf.String(); got != string(input) {
		t.Errorf("Expected incomplete sequence to pass through unchanged, got %q", got)
	}
}

func TestNonAnsiDetectionWriter_MixedSequences(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// Mix of 4-bit and 24-bit colors
	input := []byte("\x1b[31mRed\x1b[0m \x1b[38;2;100;255;100mGreen\x1b[0m \x1b[34mBlue\x1b[0m\n")
	_, err := writer.Write(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "\x1b[31mRed\x1b[0m \x1b[38;2;100;255;100m\x1b[9mGreen\x1b[29m\x1b[0m \x1b[34mBlue\x1b[0m\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}

func TestNonAnsiDetectionWriter_NonColorSequences(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// Other ANSI sequences like cursor movement should pass through
	input := []byte("\x1b[2J\x1b[H\x1b[38;2;255;100;0mColored\x1b[0m\n")
	_, err := writer.Write(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := "\x1b[2J\x1b[H\x1b[38;2;255;100;0m\x1b[9mColored\x1b[29m\x1b[0m\n"
	if got := buf.String(); got != expected {
		t.Errorf("Expected:\n%q\nGot:\n%q", expected, got)
	}
}

func TestNonAnsiDetectionWriter_BrightColors(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	// Bright/bold colors should be allowed
	input := []byte("\x1b[31;1mBright Red\x1b[0m\n")
	_, err := writer.Write(input)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if got := buf.String(); got != string(input) {
		t.Errorf("Expected bright colors to pass through unchanged, got %q", got)
	}
}

func TestNonAnsiDetectionWriter_WriterContract(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	tests := []struct {
		name     string
		input    []byte
		wantN    int // should be length of input
		wantText string
	}{
		{
			name:     "should return original length for simple text",
			input:    []byte("Hello"),
			wantN:    5,
			wantText: "Hello",
		},
		{
			name:     "should return original length with ANSI sequence",
			input:    []byte("\x1b[38;2;255;100;0mColored"),
			wantN:    24,
			wantText: "\x1b[38;2;255;100;0m\x1b[9mColored\x1b[29m",
		},
		{
			name:     "should return original length with multiple sequences",
			input:    []byte("\x1b[38;2;255;100;0mOne\x1b[0m \x1b[38;2;100;255;100mTwo"),
			wantN:    47,
			wantText: "\x1b[38;2;255;100;0m\x1b[9mOne\x1b[29m\x1b[0m \x1b[38;2;100;255;100m\x1b[9mTwo\x1b[29m",
		},
		{
			name:     "should handle incomplete sequence at end",
			input:    []byte("\x1b[38;2;255;100;0mText\x1b["),
			wantN:    23,
			wantText: "\x1b[38;2;255;100;0m\x1b[9mText\x1b[29m\x1b[",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			n, err := writer.Write(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if n != len(tt.input) {
				t.Errorf("Write() returned length %d, want %d", n, len(tt.input))
			}
			if got := buf.String(); got != tt.wantText {
				t.Errorf("Write() produced text %q, want %q", got, tt.wantText)
			}
		})
	}
}

func TestNonAnsiDetectionWriter_256Colors(t *testing.T) {
	var buf bytes.Buffer
	writer := fansischmancy.NewWriter(&buf)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "foreground 256-color",
			input:    "\x1b[38;5;82mBright Green\x1b[0m\n",
			expected: "\x1b[38;5;82m\x1b[9mBright Green\x1b[29m\x1b[0m\n",
		},
		{
			name:     "background 256-color",
			input:    "\x1b[48;5;196mRed Background\x1b[0m\n",
			expected: "\x1b[48;5;196m\x1b[9mRed Background\x1b[29m\x1b[0m\n",
		},
		{
			name:     "mixed 256 and 4-bit colors",
			input:    "\x1b[31mRed\x1b[0m \x1b[38;5;82mBright Green\x1b[0m\n",
			expected: "\x1b[31mRed\x1b[0m \x1b[38;5;82m\x1b[9mBright Green\x1b[29m\x1b[0m\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			_, err := writer.Write([]byte(tt.input))
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got := buf.String(); got != tt.expected {
				t.Errorf("Write() produced text %q, want %q", got, tt.expected)
			}
		})
	}
}
