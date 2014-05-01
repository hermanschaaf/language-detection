package language

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type transitionMap struct {
	transitions map[string]*transition
	total       int
	unique      int
}

type transition struct {
	count       int
	probability float64
}

func newTransitionMap() *transitionMap {
	m := transitionMap{}
	m.transitions = map[string]*transition{}
	return &m
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

var boundary = []byte(" \n\t\r.!<>?,;:/(){}[]0123456789~\"*&^%$#@")

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

func parse(f string) (map[string]*transitionMap, error) {
	m := map[string]*transitionMap{}
	file, err := os.Open(f)
	if err != nil {
		return m, err
	}

	s := newWordReader(file, 5) // define max word length as 2 characters

	prevTok := ""
	total, uniques := 0, 0
	for s.Scan() {
		tok := strings.ToLower(s.Text())
		if s.nodeType != wordNode {
			tok = ""
		}
		if s.nodeType == wordNode || prevTok != "" {
			total += 1
			if _, ok := m[prevTok]; !ok {
				m[prevTok] = newTransitionMap()
			}
			if _, ok := m[prevTok].transitions[tok]; !ok {
				m[prevTok].transitions[tok] = &transition{}
				m[prevTok].unique += 1
				uniques += 1
			}
			m[prevTok].total += 1
			m[prevTok].transitions[tok].count += 1
		}
		prevTok = tok
	}
	fmt.Println(total, uniques)

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func matchString(m map[string]*transitionMap, str string) (score float64) {
	s := newWordReader(strings.NewReader(str), 2)
	matches := 0
	totalWords := 0
	totalCount := 0
	prevTok := ""
	for s.Scan() {
		tok := strings.ToLower(s.Text())
		if s.nodeType == wordNode {
			totalWords += 1
			if _, ok := m[prevTok]; !ok {
				totalCount += len(m)
			} else if _, ok := m[prevTok].transitions[tok]; !ok {
				totalCount += m[prevTok].total
			} else {
				matches += m[prevTok].transitions[tok].count
				totalCount += m[prevTok].total
			}
		} else {
			prevTok = ""
		}
	}
	score = float64(matches) / float64(totalCount)
	return score
}
