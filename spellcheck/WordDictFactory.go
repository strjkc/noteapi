package spellcheck

import (
	"os"
	"path"
	"strings"
	"time"
)

type WordMapFactory struct {
	StoragePath string
	WordMaps    map[string]*WordMapCache
}

type WordMapCache struct {
	wm          *WordMap
	lastChecked time.Time
}

func NewWordMapFactory(path string) *WordMapFactory {
	wdf := WordMapFactory{StoragePath: path}
	wdf.WordMaps = map[string]*WordMapCache{}
	return &wdf
}

// TODO: caching invalidation could be a service, that runs as a go rutine independant of the other services
func (wmf *WordMapFactory) WordMap(locale string) (*WordMap, error) {
	localeFilePath := path.Join(wmf.StoragePath, locale)
	if wmc, ok := wmf.WordMaps[locale]; ok {
		currTime := time.Now()
		if currTime.After(wmc.lastChecked.Add(time.Hour)) {
			wmc.lastChecked = currTime
			info, err := os.Stat(localeFilePath)
			if err != nil {
				return nil, err
			}
			if wmc.wm.FileUpadedAt.Equal(info.ModTime()) {
				return wmc.wm, nil
			}
		} else {
			wmc.lastChecked = currTime
			return wmc.wm, nil
		}
		delete(wmf.WordMaps, locale)
	}

	data, err := os.ReadFile(localeFilePath)
	if err != nil {
		return nil, err
	}
	info, err := os.Stat(localeFilePath)
	if err != nil {
		return nil, err
	}

	wm := wmf.createWordMap(data)
	wm.FileUpadedAt = info.ModTime()

	wmc := WordMapCache{wm: wm, lastChecked: time.Now()}
	wmf.WordMaps[locale] = &wmc
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
