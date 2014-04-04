
#Overseer

Watches your specified files and executes whatever command you specify.


###How to use it

```bash
go build overseer.go
./overseer file1 file2 ... -c command goes here
# example (dogfooding)
./overseer overseer.go -c go build overseer.go
```

To run when any .language file is edited:

```bash
./overseer *.language -c ...
```

To run when any .language file is edited in any directory:

```bash
./overseer **/*.language -c ...
```

