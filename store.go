package apikey

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("key already exists")
)

type Record struct {
	ID          string
	Prefix      string
	Hash        []byte
	Fingerprint string
	CreatedAt   time.Time
	RevokedAt   *time.Time
	Disabled    bool
	Metadata    map[string]interface{}
}

func (r *Record) IsActive() bool {
	return !r.Disabled && r.RevokedAt == nil
}

type Store interface {
	Create(key string, record *Record) error
	GetByID(id string) (*Record, error)
	GetByHash(hash []byte) (*Record, error)
	Update(id string, record *Record) error
	Delete(id string) error
	List() ([]*Record, error)
	Revoke(id string) error
	Disable(id string) error
}

type MemStore struct {
	records map[string]*Record
	hashes  map[string]string
}

func NewMemStore() *MemStore {
	return &MemStore{
		records: make(map[string]*Record),
		hashes:  make(map[string]string),
	}
}

func (s *MemStore) Create(key string, record *Record) error {
	if _, ok := s.records[record.ID]; ok {
		return ErrDuplicateKey
	}

	for id := range s.hashes {
		if id == record.ID {
			return ErrDuplicateKey
		}
	}

	s.records[record.ID] = record
	record.Hash = HashKey(key)
	record.Fingerprint = Fingerprint(key)
	s.hashes[string(record.Hash)] = record.ID

	return nil
}

func (s *MemStore) GetByID(id string) (*Record, error) {
	r, ok := s.records[id]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return r, nil
}

func (s *MemStore) GetByHash(hash []byte) (*Record, error) {
	id, ok := s.hashes[string(hash)]
	if !ok {
		return nil, ErrRecordNotFound
	}
	return s.records[id], nil
}

func (s *MemStore) Update(id string, record *Record) error {
	if _, ok := s.records[id]; !ok {
		return ErrRecordNotFound
	}
	s.records[id] = record
	return nil
}

func (s *MemStore) Delete(id string) error {
	r, ok := s.records[id]
	if !ok {
		return ErrRecordNotFound
	}

	delete(s.hashes, string(r.Hash))
	delete(s.records, id)
	return nil
}

func (s *MemStore) List() ([]*Record, error) {
	result := make([]*Record, 0, len(s.records))
	for _, r := range s.records {
		result = append(result, r)
	}
	return result, nil
}

func (s *MemStore) Revoke(id string) error {
	r, ok := s.records[id]
	if !ok {
		return ErrRecordNotFound
	}
	now := time.Now()
	r.RevokedAt = &now
	return nil
}

func (s *MemStore) Disable(id string) error {
	r, ok := s.records[id]
	if !ok {
		return ErrRecordNotFound
	}
	r.Disabled = true
	return nil
}
