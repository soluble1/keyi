package memory

import (
	"context"
	"errors"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
	"web_copy/session"
)

type Store struct {
	mutex   sync.RWMutex
	cache   *cache.Cache
	expired time.Duration
}

func NewStore(expired time.Duration) *Store {
	return &Store{
		cache:   cache.New(expired, time.Second),
		expired: expired,
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess := &memorySession{
		id:   id,
		data: make(map[string]string),
	}
	s.cache.Set(sess.ID(), sess, s.expired)
	return sess, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, ok := s.cache.Get(id)
	if ok {
		return sess.(*memorySession), nil
	}
	return nil, errors.New("web：session不存在")
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.cache.Delete(id)
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, ok := s.cache.Get(id)
	if !ok {
		return errors.New("web：session不存在")
	}
	s.cache.Set(sess.(*memorySession).ID(), sess, s.expired)
	return nil
}

type memorySession struct {
	id    string
	data  map[string]string
	mutex sync.RWMutex
}

func (m *memorySession) Set(ctx context.Context, key, value string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.data[key] = value
	return nil
}

func (m *memorySession) Get(ctx context.Context, key string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return "", errors.New("未找到数据")
	}
	return val, nil
}

func (m *memorySession) ID() string {
	return m.id
}
