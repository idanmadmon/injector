package injector

import (
	"io"
	"strings"
)

func createFlie(fileSize int) io.Reader {
	return strings.NewReader(strings.Repeat("a", fileSize))
}
