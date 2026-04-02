package spellcheck

import (
	"os"
	"strings"
)

type SP struct {
	dict    map[string]struct{}
	input   string
	ch      byte
	pos     int
	readPos int
}

func NewChecker() SP {
	sp := SP{}
	sp.dict = make(map[string]struct{})
	sp.buildDict()
	return sp
}

func (s *SP) buildDict() {
	//TODO: fix path
	data, err := os.ReadFile("/home/strahinja/Work/noteapi/spellcheck/words_alpha.txt")
	if err != nil {
		panic("Cant open dict file")
	}
	dataString := string(data)
	words := strings.Split(dataString, "\r\n")
	for _, word := range words {
		s.dict[word] = struct{}{}
	}
}

func (s *SP) initState(input string) {
	s.input = input
	s.readPos = 0
	s.pos = s.readPos
	if len(s.input) == 0 {
		s.ch = byte(0)
	} else {
		s.ch = s.input[s.readPos]
	}
}

func (s *SP) CheckSpelling(input string) []string {
	s.initState(input)
	var errors []string
	s.read()
	for s.readPos < len(input) {
		if s.isChar() {
			word := s.readWord()
			if _, ok := s.dict[word]; !ok {
				errors = append(errors, word)
			}
		} else {
			s.read()
		}
	}

	return errors
}

func (s *SP) read() {
	if s.readPos < len(s.input) {
		s.ch = s.input[s.readPos]
	}
	s.pos = s.readPos
	s.readPos++
}

func (s *SP) isChar() bool {
	return s.ch >= 65 && s.ch <= 90 || s.ch >= 97 && s.ch <= 122
}

func (s *SP) readWord() string {
	pos := s.pos
	for s.isChar() && s.readPos <= len(s.input) {
		s.read()
	}
	return strings.ToLower(s.input[pos:s.pos])
}
