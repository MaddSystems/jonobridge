package persistence

import (
	"time"
)

type MemoryStateStore struct {
    data map[string]interface{}
}

func NewMemoryStateStore() *MemoryStateStore {
    return &MemoryStateStore{
        data: make(map[string]interface{}),
    }
}

func (s *MemoryStateStore) GetString(imei, key string) (string, error) {
    if v, ok := s.data[key]; ok { return v.(string), nil }
    return "", nil
}

func (s *MemoryStateStore) SetString(imei, key, value string) error {
    s.data[key] = value
    return nil
}

func (s *MemoryStateStore) GetInt64(imei, key string) (int64, error) {
    if v, ok := s.data[key]; ok { return v.(int64), nil }
    return 0, nil
}

func (s *MemoryStateStore) SetInt64(imei, key string, value int64) error {
    s.data[key] = value
    return nil
}

func (s *MemoryStateStore) GetTime(imei, key string) (time.Time, error) {
    if v, ok := s.data[key]; ok { return v.(time.Time), nil }
    return time.Time{}, nil
}

func (s *MemoryStateStore) SetTime(imei, key string, value time.Time) error {
    s.data[key] = value
    return nil
}

func (s *MemoryStateStore) GetGeofencesByGroup(groupName string) ([]Geofence, error) {
    return []Geofence{}, nil
}
