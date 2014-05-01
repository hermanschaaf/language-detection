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
	count       int
	probability float64
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

func consumeWord(data []byte, maxWordLen int) (int, []byte, error) {
	var accum []byte
	for i, b := range data {
		if bytes.IndexByte(boundary, b) > -1 || (maxWordLen > 0 && i > maxWordLen) {
			return i, accum, nil
		} else {
			accum = append(accum, b)
		}
	}
	return 0, nil, nil
}

func newWordReader(r io.Reader, maxWordLen int) *wordReader {
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
			advance, token, err = consumeWord(data, maxWordLen)
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

	m := map[string]map[string]*transition{}

	s := newWordReader(file, 2) // define max word length as 2 characters

	prevTok := ""
	total, uniques := 0, 0
	for s.Scan() {
		tok := s.Text()
		if s.nodeType == wordNode {
			total += 1
			if _, ok := m[prevTok]; !ok {
				m[prevTok] = map[string]*transition{}
			}
			if _, ok := m[prevTok][tok]; !ok {
				m[prevTok][tok] = &transition{}
				uniques += 1
			}
			m[prevTok][tok].count += 1
			prevTok = tok
		} else {
			prevTok = ""
		}
	}

	fmt.Println(total, uniques)

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return nil
}
