package local_cache

import (
	"context"
	"errors"
	cache "github.com/soluble1/mcache"
	"github.com/soluble1/mweb/session"
	"sync"
	"time"
)

type Store struct {
	mutex   sync.RWMutex
	cache   cache.Cache
	expired time.Duration
}

func NewStore(expired time.Duration) *Store {
	return &Store{
		cache:   cache.NewLocalCache(),
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
	err := s.cache.Set(ctx, sess.ID(), sess, s.expired)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	sess, err := s.cache.Get(ctx, id)
	if err == nil {
		return sess.(*memorySession), nil
	}
	return nil, err
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	err := s.cache.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sess, err := s.cache.Get(ctx, id)
	if err != nil {
		return err
	}
	err = s.cache.Set(ctx, sess.(*memorySession).ID(), sess, s.expired)
	if err != nil {
		return err
	}
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
