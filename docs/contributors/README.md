This section is dedicated for current or future ResourceD contributors.

## FAQ

**Q: What are the prerequisites to build ResourceD?**

1. Go programming language with version > 1.2.

2. [Godep](https://github.com/tools/godep).

3. [Vagrant](https://www.vagrantup.com/), to build Linux binary.


**Q: How to run the daemon?**

Below is an example on how to run ResourceD as foreground process.
```bash
RESOURCED_CONFIG_READER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-reader \
RESOURCED_CONFIG_WRITER_DIR=$GOPATH/src/github.com/resourced/resourced/tests/data/config-writer \
go run $GOPATH/src/github.com/resourced/resourced/resourced.go
```


**Q: How to run tests?**

There are a few ways to run tests:

1. On your laptop:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    godep go test ./...
    ```

2. Inside docker container:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    docker build -t resourced . && docker run -t resourced go test ./...
    ```

3. Inside Vagrant VM, docker is also pre-installed inside the VM:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    vagrant up ubuntu       # or vagrant up centos
    vagrant ssh ubuntu      # or vagrant ssh centos

    # Inside Vagrant
    export GOPATH=/go
    cd $GOPATH/src/github.com/resourced/resourced
    sudo docker build -t resourced . && sudo docker run -t resourced
    go test ./...
    ```


**Q: What is the coding style?**

Please use `go fmt` everywhere. If you use SublimeText, feel free to install `GoSublime` and `GoOracle`.


**Q: What is the general architecture?**

ResourceD has 3 components: Reader, Writer, and HTTP server.

* Reader scrapes information on your server and returns JSON data. Each reader runs in its own goroutine.

* Writer reads the JSON data, process it further, and sends it to other places. Each writer runs in its own goroutine.

* HTTP server serves the JSON data locally.
