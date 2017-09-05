package echo

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	pb "github.com/bbengfort/echo/msg"
	"github.com/bbengfort/x/stats"
	"google.golang.org/grpc"
)

func NewClient(addr, name string) (*Client, error) {
	c := new(Client)
	c.Init(addr, name)
	return c, nil
}

type Client struct {
	name     string            // host information for the server
	addr     string            // address to bind the server to
	nSent    uint64            // number of messages sent
	nRecv    uint64            // number of messages received
	nBytes   uint64            // number of bytes sent
	messages uint64            // the number of messages composed
	latency  time.Duration     // total time to send messages
	stats    *stats.Statistics // distribution of message latency
	identity string            // the identity being sent to the server
	conn     *grpc.ClientConn  // the connection to the grpc server
	stream   pb.HelloClient    // the stream to send messages on
}

func (c *Client) Init(addr, name string) {
	c.addr = addr

	// if name is empty string, set it to the hostname
	if name == "" {
		name, _ = os.Hostname()
	}
	c.name = name

	// Create an identity for the client
	// NOTE: the identity must be unique - do not rely on randomness since
	// parallel instantiation may result in the same seed!
	c.identity = fmt.Sprintf("%s-%04X", c.name, rand.Intn(0x10000))
}

func (c *Client) Connect(timeout time.Duration) (err error) {
	if c.conn, err = grpc.Dial(c.addr, grpc.WithInsecure(), grpc.WithTimeout(timeout)); err != nil {
		return WrapError("could not connect to '%s'", err, c.addr)
	}

	c.stream = pb.NewHelloClient(c.conn)
	return nil
}

func (c *Client) Close() (err error) {
	if c.conn == nil {
		return nil
	}

	if err = c.conn.Close(); err != nil {
		return WrapError("couldn't close connection", err)
	}

	c.conn = nil
	c.stream = nil
	return nil
}

func (c *Client) Shutdown() error {
	return c.Close()
}

func (c *Client) Send(msg string) error {

	req := &pb.BasicMessage{
		Sender:  c.identity,
		Message: msg,
	}

	c.nSent++
	reply, err := c.stream.Respond(context.Background(), req)
	if err != nil {
		return WrapError("could not send message", err)
	}

	c.nRecv++
	info("received: %s\n", reply.String())
	return nil
}
