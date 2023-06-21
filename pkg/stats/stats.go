package stats

import (
	"sync"
	"time"
)

type TimeEntry struct {
	TimeStamp int64
	Count     uint64
}

type Stats struct {
	RingBuff []TimeEntry
	Lock     sync.Mutex
}

func NewStats(t *[]TimeEntry) *Stats {
	stats := Stats{}
	if t == nil {
		stats.RingBuff = make([]TimeEntry, 0)
	} else {
		stats.RingBuff = *t
	}
	return &stats
}

func (s *Stats) AddEntry() uint64 {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	timeEntry := TimeEntry{
		TimeStamp: time.Now().UnixMilli(),
		Count:     1,
	}

	var currentCount uint64
	if len(s.RingBuff) == 0 {
		s.RingBuff = append(s.RingBuff, timeEntry)
		currentCount = timeEntry.Count
	} else if len(s.RingBuff) > 0 && s.RingBuff[len(s.RingBuff)-1].TimeStamp == timeEntry.TimeStamp {
		s.RingBuff[len(s.RingBuff)-1].Count++
		currentCount = s.RingBuff[len(s.RingBuff)-1].Count
	} else {
		timeEntry.Count += s.RingBuff[len(s.RingBuff)-1].Count
		s.RingBuff = append(s.RingBuff, timeEntry)
		currentCount = timeEntry.Count
	}
	return currentCount - s.RingBuff[0].Count
}

func (s *Stats) CleanupHistoricalData() {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if len(s.RingBuff) > 1 {
		if s.RingBuff[len(s.RingBuff)-1].TimeStamp-s.RingBuff[0].TimeStamp > time.Second.Milliseconds() {
			thersold := 0
			for index, entry := range s.RingBuff {
				if s.RingBuff[len(s.RingBuff)-1].TimeStamp-entry.TimeStamp < time.Second.Milliseconds() {
					thersold = index
					break
				}
			}
			s.RingBuff = s.RingBuff[thersold:]
		}
	}
}

func (s *Stats) PeriodicCleanup() {
	ticker := time.Tick(100 * time.Millisecond)
	for range ticker {
		s.CleanupHistoricalData()
	}
}
