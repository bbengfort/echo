package echo

import (
	"fmt"
	"net"
	"os"

	pb "github.com/bbengfort/echo/msg"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func NewServer(addr, name string) (*Server, error) {
	s := new(Server)
	s.Init(addr, name)
	return s, nil
}

type Server struct {
	name   string // host information for the server
	addr   string // address to bind the server to
	nSent  uint64 // number of messages sent
	nRecv  uint64 // number of messages received
	nBytes uint64 // number of bytes sent
}

func (s *Server) Init(addr, name string) {
	s.addr = addr

	// if name is empty string, set it to the hostname
	if name == "" {
		name, _ = os.Hostname()
	}
	s.name = name
}

func (s *Server) Run() error {
	sock, err := net.Listen("tcp", s.addr)
	if err != nil {
		return WrapError("could not listen on '%s'", err, s.addr)
	}
	defer sock.Close()

	info("listening for requests on %s", s.addr)

	// Create the grpc server, handler, and listen
	srv := grpc.NewServer()
	pb.RegisterHelloServer(srv, s)
	return srv.Serve(sock)
}

func (s *Server) Shutdown() error {
	return nil
}

// Respond implements the echo.HelloServer interface.
func (s *Server) Respond(ctx context.Context, in *pb.BasicMessage) (*pb.BasicMessage, error) {
	// Log that we've received the message
	s.nRecv++
	info("received: %s\n", in.String())

	// Construct the reply
	reply := &pb.BasicMessage{
		Sender:  s.name,
		Message: fmt.Sprintf("reply msg #%d", s.nRecv),
	}

	// Send the reply
	s.nSent++
	return reply, nil
}
