package diffparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// DiffParser holds the
type DiffParser struct {
	buffer *bufio.Reader
}

// New takes a Reader interface and return a new DiffParser object
func New(input io.Reader) *DiffParser {
	return &DiffParser{buffer: bufio.NewReader(input)}
}

func (d *DiffParser) read() rune {
	ch, _, err := d.buffer.ReadRune()
	if err != nil {
		return rune(0)
	}
	return ch
}

func isLineBreak(ch rune) bool {
	return ch == '\t' || ch == '\n'
}

func isWhitespace(ch rune) bool {
	return ch == ' '
}

// Parse gets every line from the buffer and parse it into the struct
func (d *DiffParser) Parse() {
	var lines [][]string
	var words []string
	word := new(bytes.Buffer)

	for {
		if ch := d.read(); ch == rune(0) {
			break
		} else if isWhitespace(ch) {
			words = append(words, word.String())
			word = new(bytes.Buffer)
			continue
		} else if isLineBreak(ch) {
			words = append(words, word.String())
			word = new(bytes.Buffer)
			lines = append(lines, words)
			words = []string{}
		} else {
			word.WriteRune(ch)
		}
	}
	diff := new(Diff)
	diff.parseFromLines(lines)
	fmt.Println(diff)
}

// Diff holds the full diff representation
type Diff struct {
	chunks []*DiffChunk
	last   *DiffChunk
}

// DiffChunk represents the differents parts of a diff chunk
type DiffChunk struct {
	filein    string
	fileout   string
	metadata  string
	chunks    []Chunk
	lastChunk *Chunk
}

// Chunk holds the lines and the changes to the code
type Chunk struct {
	lnumberIn  string
	lnumberOut string
	linesIn    [][]string
	linesOut   [][]string
}

func (d *Diff) parseFromLines(lines [][]string) {
	for _, line := range lines {
		if line[0] == "diff" && line[1] == "--git" {
			diffchunk := &DiffChunk{
				filein:  line[2],
				fileout: line[3],
			}
			d.last = diffchunk
			d.chunks = append(d.chunks, diffchunk)
		} else if line[0] == "index" {
			d.last.metadata = strings.Join(line[1:], " ")
		} else if line[0] == "+++" || line[0] == "---" {
			continue
		} else if line[0] == "@@" {
			c := &Chunk{
				lnumberOut: line[1],
				lnumberIn:  line[2],
			}
			d.last.lastChunk = c
		} else if line[0] == "-" {
			d.last.lastChunk.linesIn = append(d.last.lastChunk.linesIn, line)
		} else if line[0] == "+" {
			d.last.lastChunk.linesOut = append(d.last.lastChunk.linesOut, line)
		}
	}

	return
}
