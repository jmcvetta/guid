// Copyright 2012 Jason McVetta.  This is Free Software, released under
// an MIT-style license.  See README.md for details.

// Package guid implements a generator for roughly sorted globally unique IDs.
package guid

import (
	"fmt"
	"sync"
	"time"
)

const (
	workerIdBits       = uint64(5)
	datacenterIdBits   = uint64(5)
	maxWorkerId        = int64(-1) ^ (int64(-1) << workerIdBits)
	maxDatacenterId    = int64(-1) ^ (int64(-1) << datacenterIdBits)
	sequenceBits       = uint64(12)
	workerIdShift      = sequenceBits
	datacenterIdShift  = sequenceBits + workerIdBits
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = int64(-1) ^ (int64(-1) << sequenceBits)

	// Tue, 21 Mar 2006 20:50:14.000 GMT
	twepoch = int64(1288834974657)
)

// A GUID generator
type Generator interface {
	NextId() (int64, error) // Get the next GUID
}

type snowflake struct {
	mu         sync.Mutex
	seq        int64
	Datacenter *int64 // Datacenter ID
	Worker     *int64 // Worker ID
	LastTs     *int64 // Last timestamp
}

func milliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// NextId returns the next GUID.
func (s snowflake) NextId() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ts := milliseconds()

	if ts < *s.LastTs {
		return 0, fmt.Errorf("time is moving backwards, waiting until %d\n", *s.LastTs)
	}

	if *s.LastTs == ts {
		s.seq = (s.seq + 1) & sequenceMask
		if s.seq == 0 {
			for ts <= *s.LastTs {
				ts = milliseconds()
			}
		}
	} else {
		s.seq = 0
	}

	*s.LastTs = ts

	id := ((ts - twepoch) << timestampLeftShift) |
		(*s.Datacenter << datacenterIdShift) |
		(*s.Worker << workerIdShift) |
		s.seq

	return id, nil
}

// NewGenerator returns a Generator when configured with datacenter ID, 
// worker ID, and last timestamp.
func NewGenerator(datacenter, worker, lastts *int64) Generator {
	s := snowflake{
		Worker:     worker,
		Datacenter: datacenter,
		LastTs:     lastts,
	}
	return s
}


// NextId returns a GUID Generator without requring any special setup.  Not
// suitable for use in clustered applications.
func SimpleGenerator() Generator {
	d := int64(0)
	w := int64(0)
	l := int64(-1)
	return NewGenerator(&d, &w, &l)
}