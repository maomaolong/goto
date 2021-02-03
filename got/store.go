package got

import (
	"errors"
	"fmt"
	"strconv"

	"sync/atomic"
)

type UrlMap map[string]string

type Store struct {
	key  int64
	urls UrlMap
}

func (s *Store) Add(url string) (string, error) {
	key := s.getKey()
	_, ok := s.urls[key]
	if ok {
		return "", errors.New(key + "is exist")
	}
	s.urls[key] = url
	return s.GetShort(key)
}

func (s *Store) Get(key string) (string, error) {
	v, ok := s.urls[key]
	if !ok {
		return "", errors.New(key + "is not exist")
	}
	return v, nil
}

func (s *Store) GetShort(key string) (string, error) {
	_, ok := s.urls[key]
	if !ok {
		return "", errors.New(key + "is not exist")
	}
	return fmt.Sprintf("http://localhost:8000/%s", key), nil
}

func (s *Store) getKey() string {
	atomic.AddInt64(&s.key, 1)
	return strconv.FormatInt(s.key, 10)
}

func NewStore() *Store {
	return &Store{
		urls: make(UrlMap),
	}
}
