# NATS Module for Helium

![Codecov](https://img.shields.io/codecov/c/github/go-helium/nats.svg?style=flat-square)
[![Build Status](https://travis-ci.com/go-helium/nats.svg?branch=master)](https://travis-ci.com/go-helium/nats)
[![Report](https://goreportcard.com/badge/github.com/go-helium/nats)](https://goreportcard.com/report/github.com/go-helium/nats)
[![GitHub release](https://img.shields.io/github/release/go-helium/nats.svg)](https://github.com/go-helium/nats)
![GitHub](https://img.shields.io/github/license/go-helium/nats.svg?style=popout)

Module provides you with the following things:
- [`*nats.Conn`](https://godoc.org/github.com/nats-io/go-nats#Conn) represents a bare connection to a nats-server. It can send and receive []byte payloads
- [`stan.Conn`](https://godoc.org/github.com/nats-io/go-nats-streaming#Conn) represents a connection to the NATS Streaming subsystem. It can Publish and Subscribe to messages within the NATS Streaming cluster.

Configuration:
- yaml example
```yaml
nats:
  url: nats://<host>:<port>
  cluster_id: string
  client_id: string
  servers: [...server slice...]
  no_randomize: bool
  name: string
  verbose: bool
  pedantic: bool
  secure: bool
  allow_reconnect: bool
  max_reconnect: int
  reconnect_wait: duration
  timeout: duration
  flusher_timeout: duration
  ping_interval: duration
  max_pings_out: int
  reconnect_buf_size: int
  sub_chan_len: int
  user: string
  password: string
  token: string
```
- env example
```
NATS_URL=nats://<host>:<port>
NATS_CLUSTER_ID=string
NATS_CLIENT_ID=string
NATS_SERVERS=[...server slice...]
NATS_NO_RANDOMIZE=bool
NATS_NAME=string
NATS_VERBOSE=bool
NATS_PEDANTIC=bool
NATS_SECURE=bool
NATS_ALLOW_RECONNECT=bool
NATS_MAX_RECONNECT=int
NATS_RECONNECT_WAIT=duration
NATS_TIMEOUT=duration
NATS_FLUSHER_TIMEOUT=duration
NATS_PING_INTERVAL=duration
NATS_MAX_PINGS_OUT=int
NATS_RECONNECT_BUF_SIZE=int
NATS_SUB_CHAN_LEN=int
NATS_USER=string
NATS_PASSWORD=string
NATS_TOKEN=string
```