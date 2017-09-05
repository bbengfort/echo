# gRPC Service Template

This package is a template for designing gRPC services.

## Quickstart

Install the server/client system using `go get` as follows:

```
$ go get github.com/bbengfort/echo
```

You should then have the `echo` command installed on your system:

```
$ echo --help
```

You can run the server as follows:

```
$ echo serve
```

And send messages from the client as:

```
$ echo send "hello world"
```

Note the various arguments you can pass to both serve and send to configure the setup. Run benchmarks with the bench command:

```
$ echo bench
```

The primary comparison is between gRPC and ZMQ &mdash; the ZMQ code can be found at [github.com/bbengfort/rtreq](https://github.com/bbengfort/rtreq). 
