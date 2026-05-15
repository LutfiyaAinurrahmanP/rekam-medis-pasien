package mocks

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// InMemoryRedis adalah implementasi in-memory sederhana yang meniru behaviour
// cache.RedisStore untuk keperluan unit test. Thread-safe menggunakan sync.Mutex.
// Tidak memerlukan koneksi Redis nyata.
type InMemoryRedis struct {
	mu    sync.Mutex
	store map[string]redisEntry
	Calls []string // tracking semua operasi untuk assertion
}

type redisEntry struct {
	value     []byte
	expiresAt time.Time // zero value = tidak expired
}

// NewInMemoryRedis membuat instance baru InMemoryRedis siap pakai.
func NewInMemoryRedis() *InMemoryRedis {
	return &InMemoryRedis{
		store: make(map[string]redisEntry),
		Calls: []string{},
	}
}

// Set menyimpan value sebagai JSON dengan TTL opsional (0 = tidak expired).
// Memenuhi cache.RedisStore interface.
func (r *InMemoryRedis) Set(_ context.Context, key string, value any, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, "Set:"+key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	r.store[key] = redisEntry{value: data, expiresAt: exp}
	return nil
}

// SetRaw menyimpan raw bytes dengan TTL opsional (0 = tidak expired).
// Berguna untuk test yang perlu pre-populate data JSON yang sudah di-marshal.
func (r *InMemoryRedis) SetRaw(key string, value []byte, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, "SetRaw:"+key)

	var exp time.Time
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	r.store[key] = redisEntry{value: value, expiresAt: exp}
	return nil
}

// Get mengambil dan unmarshal value ke dest.
// Mengembalikan error jika key tidak ada atau sudah expired.
// Memenuhi cache.RedisStore interface.
func (r *InMemoryRedis) Get(_ context.Context, key string, dest any) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, "Get:"+key)

	entry, ok := r.store[key]
	if !ok {
		return errors.New("key not found")
	}
	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(r.store, key) // lazy expiry
		return errors.New("key expired")
	}
	return json.Unmarshal(entry.value, dest)
}

// Delete menghapus satu atau lebih key.
// Signature variadic (...string) sesuai dengan cache.RedisStore interface.
func (r *InMemoryRedis) Delete(_ context.Context, keys ...string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, k := range keys {
		r.Calls = append(r.Calls, "Delete:"+k)
		delete(r.store, k)
	}
	return nil
}

// DeleteByPattern menghapus semua key yang prefix-nya cocok dengan pattern.
// Karakter '*' wildcard di akhir pattern dihapus sebelum prefix matching.
// Memenuhi cache.RedisStore interface.
func (r *InMemoryRedis) DeleteByPattern(_ context.Context, pattern string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Calls = append(r.Calls, "DeleteByPattern:"+pattern)

	// Strip trailing '*' wildcard
	prefix := pattern
	if len(prefix) > 0 && prefix[len(prefix)-1] == '*' {
		prefix = prefix[:len(prefix)-1]
	}

	for k := range r.store {
		if len(k) >= len(prefix) && k[:len(prefix)] == prefix {
			delete(r.store, k)
		}
	}
	return nil
}

// ─── Helpers khusus test ──────────────────────────────────────────────────────

// Exists memeriksa apakah key ada dan tidak expired (untuk assertion di test).
func (r *InMemoryRedis) Exists(key string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	entry, ok := r.store[key]
	if !ok {
		return false
	}
	if !entry.expiresAt.IsZero() && time.Now().After(entry.expiresAt) {
		delete(r.store, key)
		return false
	}
	return true
}

// Count mengembalikan jumlah key aktif (tidak expired).
func (r *InMemoryRedis) Count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	count := 0
	now := time.Now()
	for _, entry := range r.store {
		if entry.expiresAt.IsZero() || now.Before(entry.expiresAt) {
			count++
		}
	}
	return count
}

// Reset mengosongkan seluruh store dan call log (berguna di antara test case).
func (r *InMemoryRedis) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store = make(map[string]redisEntry)
	r.Calls = []string{}
}
