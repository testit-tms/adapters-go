# How to update version

1. Update 

2. Run
```
go clean -modcache
go mod tidy
go build ./...
```

3. go test ./...
