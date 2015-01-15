Welcome (future/)contributors!

## FAQ

**Q: What are the prerequisites to build ResourceD?**

You need Go programming language with version > 1.2.


**Q: How to run tests?**

There are a few ways to run tests:

1. Inside docker container:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    docker build -t resourced . && docker run -t resourced go test ./...
    ```

2. Inside Vagrant VM, docker is also pre-installed inside the VM:
    ```bash
    cd $GOPATH/src/github.com/resourced/resourced
    vagrant up
    vagrant ssh

    # Inside Vagrant
    cd $GOPATH/src/github.com/resourced/resourced
    docker build -t resourced . && docker run -t resourced
    go test ./...
    ```


**Q: What is the coding style?**

Please use `go fmt` everywhere. If you use SublimeText, feel free to install GoSublime and GoOracle.


**Q: What is the general architecture?**

ResourceD has 3 components: Reader, writer, and HTTP server.

* Reader scrapes information in your server and returns JSON data. Each reader runs in its own goroutine.

* Writer reads the JSON data, process it further, and sends it to other places (e.g. graphite or New Relic). Each writer runs in its own goroutine.

* HTTP server serves the JSON data locally.