
#Overseer

Watches your specified files and executes whatever command you specify.


###How to use it

```bash
go build overseer.go
./overseer file1 file2 ... -c command goes here
# example (dogfooding)
./overseer overseer.go -c go build overseer.go
```


###TODO

- support regexes on files to watch
  - `*.go`
  - `**/*.go`

