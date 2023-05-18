package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger_Output(t *testing.T) {
	var buf bytes.Buffer

	logger := New()
	logger.info.SetOutput(&buf)

	logger.Info("test")

	bufStr := buf.String()
	idxByte := strings.IndexByte(bufStr, ' ')
	rest := bufStr[idxByte+1:]

	if !strings.Contains(rest, "test") {
		t.Fatalf("output %s does not contain log message", rest)
	}
}
