package media

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Protocol represents a terminal image rendering protocol.
type Protocol int

const (
	ProtocolNone   Protocol = iota
	ProtocolITerm2          // iTerm2, WezTerm
	ProtocolKitty           // Kitty terminal
)

// DetectProtocol checks environment variables to determine which
// inline image protocol the terminal supports.
func DetectProtocol() Protocol {
	term := os.Getenv("TERM_PROGRAM")
	switch term {
	case "iTerm.app", "WezTerm":
		return ProtocolITerm2
	}

	// Kitty sets TERM=xterm-kitty or TERM_PROGRAM=kitty (lowercase in some builds)
	if term == "kitty" || strings.Contains(os.Getenv("TERM"), "kitty") {
		return ProtocolKitty
	}

	return ProtocolNone
}

// RenderImageFromURL downloads an image and renders it inline.
// Returns the escape sequence string to print, or empty string on failure.
// width is in terminal columns.
func RenderImageFromURL(imageURL string, width int) string {
	if imageURL == "" {
		return ""
	}

	protocol := DetectProtocol()
	if protocol == ProtocolNone {
		return ""
	}

	data, err := fetchImage(imageURL)
	if err != nil {
		return ""
	}

	return RenderImageFromBytes(data, width, protocol)
}

// RenderImageFromBytes renders image data inline using the given protocol.
func RenderImageFromBytes(data []byte, width int, protocol Protocol) string {
	switch protocol {
	case ProtocolITerm2:
		return renderITerm2(data, width)
	case ProtocolKitty:
		return renderKitty(data, width)
	default:
		return ""
	}
}

// RenderImageFromFile reads a local file and renders it inline.
func RenderImageFromFile(path string, width int) string {
	protocol := DetectProtocol()
	if protocol == ProtocolNone {
		return ""
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}

	return RenderImageFromBytes(data, width, protocol)
}

// SupportsInlineImages returns true if the terminal supports inline image display.
func SupportsInlineImages() bool {
	return DetectProtocol() != ProtocolNone
}

func fetchImage(url string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d fetching image", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// renderITerm2 uses the iTerm2 inline image protocol.
// ESC ] 1337 ; File=[args] : base64data BEL
func renderITerm2(data []byte, width int) string {
	encoded := base64.StdEncoding.EncodeToString(data)

	var b strings.Builder
	b.WriteString("\033]1337;File=inline=1")
	fmt.Fprintf(&b, ";width=%d", width)
	b.WriteString(";preserveAspectRatio=1")
	b.WriteString(":")
	b.WriteString(encoded)
	b.WriteString("\a")

	return b.String()
}

// renderKitty uses the Kitty graphics protocol.
// ESC _G a=T,f=100,... ; base64data ESC \
// For large images, data is sent in chunks.
func renderKitty(data []byte, width int) string {
	encoded := base64.StdEncoding.EncodeToString(data)

	var b strings.Builder

	// Kitty supports chunked transfer for large images.
	// Each chunk is max 4096 bytes of base64 data.
	const chunkSize = 4096

	for i := 0; i < len(encoded); i += chunkSize {
		end := i + chunkSize
		if end > len(encoded) {
			end = len(encoded)
		}
		chunk := encoded[i:end]

		isFirst := i == 0
		isLast := end >= len(encoded)

		b.WriteString("\033_G")
		if isFirst {
			// First chunk: include action and format
			fmt.Fprintf(&b, "a=T,f=100,c=%d", width)
		}

		switch {
		case isFirst && isLast:
			// Single chunk, no more data
			b.WriteString(";")
		case isFirst:
			// First of multiple chunks
			b.WriteString(",m=1;")
		case isLast:
			// Last chunk
			b.WriteString("m=0;")
		default:
			// Middle chunk
			b.WriteString("m=1;")
		}

		b.WriteString(chunk)
		b.WriteString("\033\\")
	}

	return b.String()
}
