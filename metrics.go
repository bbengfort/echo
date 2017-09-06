package echo

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

//===========================================================================
// Server-Side Metrics
//===========================================================================

// Metrics tracks the measurable statistics of the system over time from the
// perspective of the local replica. Many stats are simply counters, other
// statistics perform online computations of the distribution of values.
type Metrics struct {
	sync.RWMutex
	started  time.Time         // The time of the first client message
	finished time.Time         // The time of the last client message
	accesses map[string]uint64 // The number of messages per-client recv by the server
}

// Init the metrics
func (m *Metrics) Init() {
	m.accesses = make(map[string]uint64)
}

// Accesses returns the total number of accesses to the replica.
func (m *Metrics) Accesses() uint64 {
	m.RLock()
	defer m.RUnlock()

	var total uint64
	for _, count := range m.accesses {
		total += count
	}
	return total
}

// Increment the access metrics and set the started time.
func (m *Metrics) Increment(client string) {
	m.Lock()
	defer m.Unlock()

	if m.started.IsZero() {
		m.started = time.Now()
	}

	m.accesses[client]++
}

// Complete an access and set the finished time.
func (m *Metrics) Complete() {
	m.Lock()
	defer m.Unlock()

	m.finished = time.Now()
}

// Duration computes the amount of time during which accesses were received.
func (m *Metrics) Duration() time.Duration {
	m.RLock()
	defer m.RUnlock()
	return m.finished.Sub(m.started)
}

// Throughput computes the number of messages per second.
func (m *Metrics) Throughput() (throughput float64) {
	m.RLock()
	defer m.RUnlock()

	duration := m.Duration()
	accesses := m.Accesses()

	if accesses > 0 && duration > 0 {
		throughput = float64(accesses) / duration.Seconds()
	}

	return throughput
}

// NClients returns the number of clients accessing the replica
func (m *Metrics) NClients() uint64 {
	m.RLock()
	defer m.RUnlock()

	var n uint64
	for _ = range m.accesses {
		n++
	}

	return n
}

// ClientMean returns the average number of accesses per client.
func (m *Metrics) ClientMean() float64 {
	m.RLock()
	defer m.RUnlock()

	var n uint64
	var s uint64
	for _, count := range m.accesses {
		n++
		s += count
	}

	if n > 0 {
		return float64(s) / float64(n)
	}

	return 0.0
}

// Serialize the data structure to a map
func (m *Metrics) Serialize(extra map[string]interface{}) map[string]interface{} {
	m.RLock()
	defer m.RUnlock()

	data := make(map[string]interface{})
	data["clients"] = m.NClients()
	data["accesses"] = m.Accesses()
	data["mean"] = m.ClientMean()
	data["duration"] = m.Duration().String()
	data["throughput"] = m.Throughput()

	for key, val := range extra {
		data[key] = val
	}

	return data
}

// String returns a quick summary of the access metrics
func (m *Metrics) String() string {
	m.RLock()
	defer m.RUnlock()

	return fmt.Sprintf(
		"%d accesses by %d clients in %s -- %0.4f accesses/second",
		m.Accesses(), m.NClients(), m.Duration(), m.Throughput(),
	)
}

// Append another metrics' data to the current metrics
func (m *Metrics) Append(o *Metrics) {
	m.Lock()
	o.RLock()
	defer m.Unlock()
	defer o.RUnlock()

	// Increment the counts
	for client, count := range o.accesses {
		m.accesses[client] += count
	}

	// If the other started time is earlier, set it as started
	if !o.started.IsZero() && (m.started.IsZero() || o.started.Before(m.started)) {
		m.started = o.started
	}

	// If the other finished time is later, set it as finished
	if !o.finished.IsZero() && (m.finished.IsZero() || o.finished.After(m.finished)) {
		m.finished = o.finished
	}
}

// Write the metrics to the path, appending the JSON as a line to the file.
func (m *Metrics) Write(path string, extra map[string]interface{}) error {
	// Don't do anything if no path is given
	if path == "" {
		return nil
	}

	// Open the file for appending, creating it if necessary
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Marshal the JSON in one line without indents
	data, err := json.Marshal(m.Serialize(extra))
	if err != nil {
		return err
	}

	// Append a newline to the data
	data = append(data, byte('\n'))

	// Append the data to the file
	_, err = f.Write(data)
	return err
}
