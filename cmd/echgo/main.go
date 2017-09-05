package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bbengfort/echo"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

//===========================================================================
// Main Method
//===========================================================================

func main() {

	// Load the .env file if it exists
	godotenv.Load()

	// Instantiate the command line application
	app := cli.NewApp()
	app.Name = "echgo"
	app.Version = "0.1"
	app.Usage = "run gRPC echo server and client"

	// Define commands available to the application
	app.Commands = []cli.Command{
		{
			Name:     "serve",
			Usage:    "run the echo server",
			Category: "server",
			Action:   serve,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "address to bind the server to",
					Value: ":4157",
				},
				cli.StringFlag{
					Name:  "n, name",
					Usage: "name to identify the server (default is hostname)",
				},
				cli.StringFlag{
					Name:  "u, uptime",
					Usage: "pass a parsable duration to shut the server down after",
				},
				cli.UintFlag{
					Name:  "verbosity",
					Usage: "set log level from 0-4, lower is more verbose",
					Value: 3,
				},
			},
		},
		{
			Name:     "send",
			Usage:    "send a message to the server",
			Category: "client",
			Action:   send,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "address to connect to the server on",
					Value: "localhost:4157",
				},
				cli.StringFlag{
					Name:  "n, name",
					Usage: "name to identify the client (default is hostname)",
				},
				cli.StringFlag{
					Name:  "t, timeout",
					Usage: "recv timeout for each message",
					Value: "5s",
				},
				cli.IntFlag{
					Name:  "r, retries",
					Usage: "number of retries before quitting",
					Value: 3,
				},
			},
		},
		{
			Name:     "bench",
			Usage:    "run throughput benchmarks",
			Category: "client",
			Action:   bench,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "a, addr",
					Usage: "address to connect to the server on",
					Value: "localhost:4157",
				},
				cli.StringFlag{
					Name:  "n, name",
					Usage: "name to identify the server (default is hostname)",
				},
				cli.StringFlag{
					Name:  "d, duration",
					Usage: "parsable duration of the benchmark",
					Value: "30s",
				},
				cli.StringFlag{
					Name:  "t, timeout",
					Usage: "recv timeout for each message",
					Value: "5s",
				},
				cli.IntFlag{
					Name:  "r, retries",
					Usage: "number of retries before quitting",
					Value: 3,
				},
				cli.IntFlag{
					Name:  "c, clients",
					Usage: "extra information: number of clients",
				},
				cli.StringFlag{
					Name:  "o, results",
					Usage: "path to write the results to",
					Value: "results.json",
				},
				cli.Int64Flag{
					Name:  "s, seed",
					Usage: "specify random seed for the process",
					Value: time.Now().Unix(),
				},
				cli.UintFlag{
					Name:  "verbosity",
					Usage: "set log level from 0-4, lower is more verbose",
					Value: 3,
				},
			},
		},
	}

	// Run the CLI program
	app.Run(os.Args)
}

//===========================================================================
// Server Commands
//===========================================================================

func exit(msg string, err error, a ...interface{}) error {
	if msg != "" {
		msg = fmt.Sprintf(msg, a...)
		msg += ": %s"
	} else {
		msg = "fatal error: %s"
	}
	return cli.NewExitError(fmt.Sprintf(msg, err), 1)
}

func serve(c *cli.Context) error {
	// Set the debug log level
	verbose := c.Uint("verbosity")
	echo.SetLogLevel(uint8(verbose))

	// Create the server
	server, err := echo.NewServer(c.String("addr"), c.String("name"))
	if err != nil {
		return exit("could not initialize server", err)
	}

	// Defer the shutdown
	defer server.Shutdown()

	// If uptime is specified, set a fixed duration for the server to run.
	if uptime := c.String("uptime"); uptime != "" {
		d, err := time.ParseDuration(uptime)
		if err != nil {
			return exit("could not parse uptime", err)
		}

		time.AfterFunc(d, func() {
			server.Shutdown()
			os.Exit(0)
		})
	}

	// Run the network server and broadcast clients
	if err := server.Run(); err != nil {
		return exit("could not run server", err)
	}
	return nil
}

//===========================================================================
// Client Commands
//===========================================================================

func send(c *cli.Context) error {
	client, err := echo.NewClient(c.String("addr"), c.String("name"))
	if err != nil {
		return exit("could not create client", err)
	}
	defer client.Shutdown()

	var timeout time.Duration
	if timeout, err = time.ParseDuration(c.String("timeout")); err != nil {
		return exit("", err)
	}

	if err = client.Connect(timeout); err != nil {
		return exit("", err)
	}

	for _, msg := range c.Args() {
		if err := client.Send(msg); err != nil {
			exit("", err)
		}
	}

	return client.Close()
}

func bench(c *cli.Context) error {

	// Set the debug log level
	verbose := c.Uint("verbosity")
	echo.SetLogLevel(uint8(verbose))

	// Set the random seed
	rand.Seed(c.Int64("seed"))

	client, err := echo.NewClient(c.String("addr"), c.String("name"))
	if err != nil {
		return exit("could not create client", err)
	}
	defer client.Shutdown()

	var duration time.Duration
	if duration, err = time.ParseDuration(c.String("duration")); err != nil {
		return exit("", err)
	}

	var timeout time.Duration
	if timeout, err = time.ParseDuration(c.String("timeout")); err != nil {
		return exit("", err)
	}

	if err = client.Connect(timeout); err != nil {
		return exit("", err)
	}
	defer client.Close()

	nClients := c.Int("clients")
	// retries := c.Int("retries")
	results := c.String("results")

	return client.Benchmark(duration, results, nClients)
}
