package master

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"sync"
)

type UrlMap map[string]string

type record struct {
	Key, Url string
}

type Store struct {
	saveChan chan record
	file     *os.File
	rw       sync.RWMutex
	urls     UrlMap
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

func (s *Store) Add(key, url *string) error {
	s.rw.Lock()
	defer s.rw.Unlock()
	_, ok := s.urls[*key]
	if ok {
		return errors.New(*key + "is exist")
	}
	s.urls[*key] = *url
	s.saveChan <- record{*key, *url}
	return nil
}

func (s *Store) Get(key, url *string) error {
	s.rw.RLock()
	defer s.rw.RUnlock()
	v, ok := s.urls[*key]
	if !ok {
		return errors.New(*key + "is not exist")
	}
	*url = v
	return nil
}

func (s *Store) GetUrls() UrlMap {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return s.urls
}

func (s *Store) set(key, url string) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.urls[key] = url
}

func (s *Store) load() error {
	decoder := json.NewDecoder(s.file)
	var err error
	for err == nil {
		var r record
		if err = decoder.Decode(&r); err == nil {
			s.set(r.Key, r.Url)
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
