package got

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type UrlMap map[string]string

type record struct {
	Key, Url string
}

type Store struct {
	saveChan chan record
	file     *os.File
	rw       sync.RWMutex
	key      int64
	urls     UrlMap
}

func (s *Store) Set(key, url string) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.urls[key] = url
}

func (s *Store) Add(url string) (string, error) {
	s.rw.Lock()
	defer s.rw.Unlock()
	key := s.createKey()
	_, ok := s.urls[key]
	if ok {
		return "", errors.New(key + "is exist")
	}
	s.urls[key] = url
	s.saveChan <- record{key, url}
	return fmt.Sprintf("http://localhost:8000/%s", key), nil
}

func (s *Store) Get(key string) (string, error) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	v, ok := s.urls[key]
	if !ok {
		return "", errors.New(key + "is not exist")
	}
	return v, nil
}

func (s *Store) createKey() string {
	atomic.AddInt64(&s.key, 1)
	return strconv.FormatInt(s.key, 10)
}

func (s *Store) load() error {
	decoder := json.NewDecoder(s.file)
	var err error
	for err == nil {
		var r record
		if err = decoder.Decode(&r); err == nil {
			var err1 error
			s.key, err1 = strconv.ParseInt(r.Key, 10, 64)
			if err1 != nil {
				log.Println(err1.Error())
			}
			s.Set(r.Key, r.Url)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *Store) save() {
	for {
		r, ok := <-s.saveChan
		if ok {
			encoder := json.NewEncoder(s.file)
			err := encoder.Encode(r)
			if err != nil {
				log.Println(err.Error())
			}
		} else {
			break
		}
	}
}

func NewStore() *Store {
	store := &Store{
		saveChan: make(chan record, 100),
		urls:     make(UrlMap),
	}
	f, err := os.OpenFile("store.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	store.file = f
	err = store.load()
	if err != nil {
		log.Fatal(err.Error())
	}
	go store.save()
	return store
}
