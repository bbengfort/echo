package echo

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/bbengfort/x/stats"
)

// Benchmark the throughput in terms of messages per second to the zmqnet.
func (c *Client) Benchmark(duration time.Duration, results string, nClients int) error {

	// Initialize the client
	c.messages = 0
	c.latency = 0
	c.nSent = 0
	c.nRecv = 0
	c.nBytes = 0
	c.stats = new(stats.Statistics)

	// Initialize the results
	extra := make(map[string]interface{})
	extra["n_clients"] = nClients
	extra["name"] = c.identity

	// Initialize channels
	timer := time.NewTimer(duration)
	echan := make(chan error, 1)
	done := make(chan bool, 1)
	status("starting benchmark for %s", duration)

	// Send the first access
	go c.Access(done, echan)

	// Continue until the timer is complete
	for {
		select {
		case <-timer.C:
			// Benchmarking complete
			return c.Results(results, extra)
		case err := <-echan:
			// Something went wrong
			return err
		case <-done:
			go c.Access(done, echan)
		}
	}

}

// Access sends a request to the server and waits for a response, measuring
// the latency of the message send to get throughput benchmarks.
func (c *Client) Access(done chan<- bool, echan chan<- error) {
	// Prepare the send
	message := fmt.Sprintf("msg %d at %s", c.messages+1, time.Now())
	start := time.Now()

	// Send the request
	if err := c.Send(message); err != nil {
		echan <- err
		return
	}

	// Compute the throughput
	latency := time.Since(start)
	c.messages++
	c.latency += latency
	c.stats.Update(float64(latency))

	// Signal done
	done <- true
}

// Results saves the throughput to disk
func (c *Client) Results(path string, data map[string]interface{}) error {
	debug("writing results to %s", path)
	data["messages"] = c.messages
	data["latency (nsec)"] = c.latency.Nanoseconds()
	data["throughput (msg/sec)"] = float64(c.messages) / c.latency.Seconds()
	data["latency distribution"] = c.stats.Serialize()
	status("%d messages in %0.3f seconds - %0.3f msg/sec", c.messages, c.latency.Seconds(), data["throughput (msg/sec)"])
	return appendJSON(path, data)
}

// Helper function to append json data as a one line string to the end of a
// results file without deleting the previous contents in it.
func appendJSON(path string, val interface{}) error {
	// Open the file for appending, creating it if necessary
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Marshal the JSON in one line without indents
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	// Append a newline to the data
	data = append(data, byte('\n'))

	// Append the data to the file
	_, err = f.Write(data)
	return err
}
