Set your GO env variables

```
GOPATH=~/go
PATH=$PATH:$GOPATH/bin
```

Prepare your workspace

```
mkdir -p ~/go/src/github.com/grokify
```

Clone the project

```
cd ~/go/src/github.com/grokify
clone git@github.com:grokify/chathooks.git
cd chathooks
```

Install godep (dependency manager)

```
go get github.com/tools/godep
```

Download all dependencies

```
godep restore
```

Run the project then visit http://localhost:3000/

```
go run main.go
```
