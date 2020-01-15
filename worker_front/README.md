# Worker Front

Task runner of woker nodes. Excutes docker container of the function requested by the load balancer.

### What does it do

- Provide RESTFul API to run a function in the node
- Manage docker images in the node (auto remove with the lifecycle)
- Build image(vm environment) automatically
	- Download the target function source code (from S3)
	- Attach environment by function's demands (python3.7, NodeJS 12, Java 8)
	- Build docker image with the function source code

### How to run

```bash
go run *.go
```

Run as a daemon (with `port=8222`):

```bash
go run *.go 8222 -d
```

Kill daemon :

```bash
kill `cat daemon.pid`
```

### How to install

```bash
go build
```

or

```bash
go get github.com/docker/docker@master
go get github.com/mholt/archiver
go get github.com/pierrec/lz4@v3
go get github.com/sevlyar/go-daemon
```
