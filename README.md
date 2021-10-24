# server-scratch

A small scratch to mess around with various server implementations. For right now, we mess around with the stdlib (AKA net/http), [evio](https://github.com/tidwall/evio), and [gnet](https://github.com/panjf5000/gnet).

## Goals

Create functional stdlib compatible server implementations to check out HTTP.

Evio and Gnet are essentially almost the exact same API, with slightly different implementation details, so it was easy to create something that would read the incoming data frame and then parse it into the http request struct.

This whole thing was wrapped so that I can play with various implementations and see how they work out.