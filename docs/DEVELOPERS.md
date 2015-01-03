## Running ResourceD tests inside Docker container

#### 1. Building Docker image

```
cd $GOPATH/src/github.com/resourced/resourced
docker build -t resourced .
```

#### 2. Running tests inside Docker container

```
cd $GOPATH/src/github.com/resourced/resourced
docker run -i -t resourced /bin/bash
go test ./...
```
