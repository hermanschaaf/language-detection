package language

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

type transition struct {
	probability float64
	target      string
}

const (
	wordNode = iota
	boundaryNode
)

type wordReader struct {
	*bufio.Scanner
	nodeType int
}

func (w *wordReader) tokenType() int {
	return w.nodeType
}

var boundary = []byte{' ', '\n', '\t', '\r', '.', ',', ';'}

func consumeBoundary(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if bytes.IndexByte(boundary, b) == -1 {
			return i, accum, nil
		} else {
			accum = append(accum, b)
		}
	}
	return 0, nil, nil
}

func consumeWord(data []byte) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if bytes.IndexByte(boundary, b) > -1 {
			return i, accum, nil
		} else {
			accum = append(accum, b)
		}
	}
	return 0, nil, nil
}

func newWordReader(r io.Reader) *wordReader {
	s := bufio.NewScanner(r)
	rdr := &wordReader{
		Scanner: s,
	}
	// splitFunc defines how we split our tokens
	splitFunc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if bytes.IndexByte(boundary, data[0]) > -1 {
			advance, token, err = consumeBoundary(data)
			rdr.nodeType = boundaryNode
		} else {
			advance, token, err = consumeWord(data)
			rdr.nodeType = wordNode
		}
		return
	}
	s.Split(splitFunc)
	return rdr
}

func parse(f string) error {
	file, err := os.Open(f)
	if err != nil {
		return err
	}

	// m := map[string]transition{}

	s := newWordReader(file)
	for s.Scan() {
		tok := s.Text()
		if s.nodeType == wordNode {
			fmt.Println(tok, s.nodeType)
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
