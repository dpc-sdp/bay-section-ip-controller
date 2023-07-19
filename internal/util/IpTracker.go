package util

import (
	"sync"
	"time"
)

type IPTracker struct {
	records map[string]Record
	mu      sync.Mutex
}

type Record struct {
	Count    int       `json:"count"`
	LastSeen time.Time `json:"lastSeen"`
}

func NewIPTracker() *IPTracker {
	return &IPTracker{
		records: make(map[string]Record),
	}
}

func (t *IPTracker) TrackIP(ip string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	record, exists := t.records[ip]
	if exists {
		record.Count++
		record.LastSeen = time.Now()
	} else {
		record = Record{
			Count:    1,
			LastSeen: time.Now(),
		}
	}

	t.records[ip] = record
}

func (t *IPTracker) GetAll() map[string]Record {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.records
}

func (t *IPTracker) BackoffIPs() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	var l []string

	for ip, record := range t.records {
		since := time.Since(record.LastSeen)
		// Keep an IP in the block list for up to 1 hour if
		// it is a repeat abuser.
		if since < 60*time.Minute && record.Count > 5 {
			l = append(l, ip)
		}
	}

	return l
}
