package spellcheck

import (
	"os"
	"path"
	"strings"
)

type WordMapFactory struct {
	StoragePath string
	WordMaps    map[string]*WordMap
}

func NewWordMapFactory(path string) *WordMapFactory {
	wdf := WordMapFactory{StoragePath: path}
	wdf.WordMaps = map[string]*WordMap{}
	return &wdf
}

func (wmf *WordMapFactory) WordMap(locale string) (*WordMap, error) {
	if wmf, ok := wmf.WordMaps[locale]; ok {
		return wmf, nil
	}

	data, err := os.ReadFile(path.Join(wmf.StoragePath, locale))
	if err != nil {
		return nil, err
	}
	wm := wmf.createWordMap(data)
	wmf.WordMaps[locale] = wm
	return wm, nil
}

func (wmf *WordMapFactory) createWordMap(data []byte) *WordMap {
	wm := WordMap{}
	wm.ValidWords = make(map[string]struct{}, 0)
	dataString := string(data)
	words := strings.Split(dataString, "\r\n")
	for _, word := range words {
		wm.ValidWords[word] = struct{}{}
	}
	return &wm
}
