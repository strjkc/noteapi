package spellcheck

import (
	"os"
	"testing"
	"time"
)

func TestCached(t *testing.T) {
	t.Run("WM Is Retreived From cache, file check when file not changed", func(t *testing.T) {
		info, err := os.Stat("/Users/strahinja/Work/noteapi/spellcheck/eng")
		if err != nil {
			t.Fatal("Failed on setup")
			return
		}
		modTime := info.ModTime()
		expectedObj := WordMap{ValidWords: map[string]struct{}{"word": {}}, FileUpadedAt: modTime}
		cache := WordMapCache{wm: &expectedObj, lastChecked: time.Now().Add(time.Hour * -2)}
		wmf := NewWordMapFactory("/Users/strahinja/Work/noteapi/spellcheck/")
		wmf.WordMaps["eng"] = &cache

		wm, err := wmf.WordMap("eng")
		if err != nil {
			t.Fatal(err)
			return
		}
		if wm != &expectedObj {
			t.Fatal("Object not retreived from the cache")
			return
		}
	})

	t.Run("WM Is Retreived From cache, cache refresh time not expired", func(t *testing.T) {
		info, err := os.Stat("/Users/strahinja/Work/noteapi/spellcheck/eng")
		if err != nil {
			t.Fatal("Failed on setup")
			return
		}
		modTime := info.ModTime()
		expectedObj := WordMap{ValidWords: map[string]struct{}{"word": {}}, FileUpadedAt: modTime}
		cache := WordMapCache{wm: &expectedObj, lastChecked: time.Now().Add(time.Minute * -58)}
		wmf := NewWordMapFactory("/Users/strahinja/Work/noteapi/spellcheck/")
		wmf.WordMaps["eng"] = &cache

		wm, err := wmf.WordMap("eng")
		if err != nil {
			t.Fatal(err)
			return
		}
		if wm != &expectedObj {
			t.Fatal("Object not retreived from the cache")
			return
		}
	})

	t.Run("WM not cached, created new", func(t *testing.T) {
		wmf := NewWordMapFactory("/Users/strahinja/Work/noteapi/spellcheck/")

		wm, err := wmf.WordMap("eng")
		if err != nil {
			t.Fatal(err)
			return
		}
		if wm == nil {
			t.Fatal("Object not created")
			return
		}

		if _, ok := wmf.WordMaps["eng"]; !ok {
			t.Fatal("Object not cached")
			return
		}
	})

	t.Run("WM Is Created New, after file chcek", func(t *testing.T) {
		wmf := NewWordMapFactory("/Users/strahinja/Work/noteapi/spellcheck/")
		wm, err := wmf.WordMap("eng")
		if err != nil {
			t.Fatal(err)
			return
		}
		wmc := wmf.WordMaps["eng"]
		wmc.lastChecked = time.Now().Add(time.Hour * -2)
		currTime := time.Now()
		os.Chtimes("/Users/strahinja/Work/noteapi/spellcheck/eng", currTime, currTime)
		wm2, err := wmf.WordMap("eng")
		if err != nil {
			t.Fatal(err)
			return
		}

		if wm2 == wm {
			t.Fatal("Object not created again")
			return
		}
	})
}
