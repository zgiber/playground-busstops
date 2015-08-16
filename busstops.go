package busstops

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type FeatureCollection struct {
	Features []json.RawMessage `json:"features"`
}

type index map[string]map[string]interface{}

type Store struct {
	l       sync.RWMutex
	i       index
	indices []string
}

func NewStore(indices ...string) *Store {
	return &Store{
		i:       map[string]map[string]interface{}{},
		indices: indices,
	}
}

func (s *Store) Add(data map[string]interface{}) {
	s.l.Lock()
	defer s.l.Unlock()
	parseIndex(data, s.i, s.indices...)
}

// TODO: obviously, a Remove / Del method if the dataset changes

func (s *Store) Get(key, value string) (interface{}, bool) {
	s.l.RLock()
	defer s.l.RUnlock()
	if values, ok := s.i[key]; ok {
		if v, ok := values[value]; ok {
			return v, true
		}
	} else {
		// invalid key
	}

	return nil, false
}

func parseIndex(data map[string]interface{}, i index, indices ...string) {
	flattened := map[string]interface{}{} // init the flattened entry outside the loop to reuse
	flatten(data, flattened, "")          // got one entry flat
	for _, index := range indices {
		if indexedValue, ok := flattened[index]; ok {
			if _, ok := i[index]; !ok {
				i[index] = map[string]interface{}{}
			}

			i[index][fmt.Sprint(indexedValue)] = data
		}
	}
}

func flatten(src, dst map[string]interface{}, prefix string) {
	for k, v := range src {
		switch t := v.(type) {
		case map[string]interface{}:
			flatten(t, dst, k)
		default:
			if len(prefix) > 0 {
				k = strings.Join([]string{prefix, "_", k}, "")
			}
			dst[k] = v
		}
	}
}
