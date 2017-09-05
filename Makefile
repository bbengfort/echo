# Shell to use with Make
SHELL := /bin/bash

# Export targets not associated with files.
.PHONY: protobuf

# Compile protocol buffers
protobuf:
	@echo "compiling protocol buffers"
	@protoc -I msg/ msg/*.proto --go_out=plugins=grpc:msg/
