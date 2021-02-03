package got

import (
	"errors"
	"fmt"
	"strconv"
)

type UrlMap map[string]string

type Store struct {
	key  int
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
	s.key += 1
	return strconv.Itoa(s.key)
}

func NewStore() *Store {
	return &Store{
		urls: make(UrlMap),
	}
}
