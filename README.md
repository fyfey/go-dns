# go-dns

A simple dns server written in Go.

This project is purely to learn how DNS works at a protocol level, and to play with Go.

### Run
```
go run .
```
Run a query against the server:
```
dig fyfe.io @localhost -p 1234
```

### Run unit tests
```
go test ./...
```
