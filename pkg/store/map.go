package store

import (
	"context"
	"sync"
	"time"
)

type record struct {
	value interface{}
	last  int64
}

type Store struct {
	records map[interface{}]*record
	mux     sync.Mutex
	max     int64
	limit   int
}

// New creates a store for time to live records.
// `ctx` controls livecycle.
// `size` is the initial size of the store.
// `ttl` is time to live.
// GC period is 5 seconds.
// Default record limit is 50.
func New(ctx context.Context, size int, ttl time.Duration) *Store {
	s := &Store{
		records: make(map[interface{}]*record, size),
		max:     int64(ttl.Seconds()),
		limit:   50,
	}
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.gc()
			}
		}
	}()
	return s
}

func (s *Store) gc() {
	s.mux.Lock()
	defer s.mux.Unlock()
	now := time.Now().Unix()

	var oldest interface{}
	oldestAge := int64(0)
	defer func() {
		if len(s.records) > s.limit {
			// log.Println("[limit] deleting record", oldest)
			delete(s.records, oldest)
		}
	}()

	for k, v := range s.records {
		age := now - v.last
		if age > s.max {
			// log.Println("[ttl] deleting record", k)
			delete(s.records, k)
		} else if age > oldestAge {
			oldestAge = age
			oldest = k
		}
	}
}

func (s *Store) Keys() []interface{} {
	keys := make([]interface{}, 0)
	for k := range s.records {
		keys = append(keys, k)
	}
	return keys
}

func (s *Store) Set(k interface{}, v interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	r, ok := s.records[k]
	if !ok {
		r = &record{value: v}
		s.records[k] = r
	} else {
		r.value = v
	}

	r.last = time.Now().Unix()
}

func (s *Store) Get(k interface{}) interface{} {
	s.mux.Lock()
	defer s.mux.Unlock()

	if it, ok := s.records[k]; ok {
		it.last = time.Now().Unix()
		return it.value
	}

	return nil
}

func (s *Store) Delete(k interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.records, k)
}

func (s *Store) Update(k interface{}, fn func(value interface{}, found bool) interface{}) {
	s.mux.Lock()
	defer s.mux.Unlock()

	r, found := s.records[k]
	var v interface{}

	if found {
		v = r.value
	}

	// invoke callback.
	v = fn(v, found)

	if !found {
		r = &record{}
		s.records[k] = r
	}

	r.value = v
	r.last = time.Now().Unix()
}
