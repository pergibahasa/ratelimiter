## How To Use

1. For testing, use file example.go and then compile to executable file

2. run the compiled executable file

3. open other terminal prompt, run this command

```bash
echo "GET http://localhost:8888/" | vegeta attack -duration=10s -rate=100 | tee results.bin | vegeta report
```

### Note

`required : go 1.18 or newer`

```bash
go get -u github.com/pergibahasa/ratelimiter
```
