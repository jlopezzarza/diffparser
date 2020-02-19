package diffparser

import (
	"bufio"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	f, err := os.Open("examples/correct_example.diff")
	if err != nil {
		t.Errorf("Diff example not found: %v", err)
	}
	filebuffer := bufio.NewReader(f)
	dp := New(filebuffer)
	dp.Parse()
}
