# KIT library

## Description

`Kit` is a library providing with implementations of infrastructure components

The main idea is to have infrastructure code reusable across all the services

## Installation

* make sure you have `make` installed
* make sure you have `golangci-lint` installed

###### Local Build
````
# build
make build

# build and generate all artifacts
make artifacts
````

## Import

It should be imported as a `module`

Note there is a trick of how to import a module from private repo to your project

Draw attention to `replace` section in the snippet below 

````
require (
	github.com/exluap/kit [kit-version]
)

replace (
	github.com/exluap/kit => github.com/exluap/kit.git [kit-version]
```` 
 
## Tests

Run tests with a command
````
make test
````

Run integration tests with a command
````
make test-integration
````

## Package summary

|name|description|
|----|-----------|
|./bpm|BPM-engine related staff, implementation of Zeebe access and helpers|
|./cache|distributed cache access, Redis|
|./common|base common objects|
|./config|utilities to work with configuration files|
|./context|utilities to context object|
|./cron|wrapper around cron library to schedule tasks|
|./grpc|grpc-related staff|
|./http|http-related staff|
|./kv|KV-store access and utilities|
|./log|logger implementation|
|./queue|message brokers' access and utilities, currently NATS & NATS streaming|
|./search|index search, Elastic Search|
|./service|utilities common for all services, like coordination cluster mechanism|
|./db|databases-related staff|




