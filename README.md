# Go package provides outgoing queue for channel to use in select case

[![Go Reference](https://pkg.go.dev/badge/github.com/powerman/chanq.svg)](https://pkg.go.dev/github.com/powerman/chanq)
[![CI/CD](https://github.com/powerman/chanq/workflows/CI/CD/badge.svg?event=push)](https://github.com/powerman/chanq/actions?query=workflow%3ACI%2FCD)
[![Coverage Status](https://coveralls.io/repos/github/powerman/chanq/badge.svg?branch=master)](https://coveralls.io/github/powerman/chanq?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/powerman/chanq)](https://goreportcard.com/report/github.com/powerman/chanq)
[![Release](https://img.shields.io/github/v/release/powerman/chanq)](https://github.com/powerman/chanq/releases/latest)

## Performance

About 1.5-1.7 times slower than sending to a blocking channel.

## Example

```go
out := make(chan []byte) // Usually not buffered, with blocking send.
q := chanq.NewQueue(out)
q.Enqueue([]byte(`one`))
q.Enqueue([]byte(`two`))
for {
    select {
    case data := <-in: // E.g.: forward from in to out without blocking.
        q.Enqueue(data)
    case q.C <- q.Elem: // Works only when queue is not empty.
        q.Dequeue()
    }
}
```
