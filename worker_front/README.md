# Worker Front

Task runner of woker nodes. Excutes docker container of the function requested by the load balancer.

### What does it do

- Provide RESTFul API to run a function in the node
- Manage docker images in the node (auto remove with the lifecycle)
- Build image(vm environment) automatically
	- Download the target function source code (from S3)
	- Attach environment by function's demands (python3.7, NodeJS 12, Java 8)
	- Build docker image with the function source code


### How to install

Please use golang >= 1.13 since we use `gomod` for dependency control.  
After update, run command below:

```bash
go build
```

Or if you are not able to use `gomod`, get 3rd party libraries manually.

```bash
go get -u github.com/docker/docker@master
go get -u github.com/mholt/archiver
go get -u github.com/pierrec/lz4@v3
go get -u github.com/sevlyar/go-daemon
```


### How to run

We support two arguments for worker front.
First argument is the port, and second argument is the option for daemon. 
Default port number is `8222`.

```bash
$ go run *.go start
$ go run *.go start 8222
$ go run *.go start 8222 -d
```

If you use daemonize option,
then you can see two files (`daemon.log`, `daemon.pid`) on the same directory.
To kill the daemon, run command like below.

```bash
$ go run *.go stop
```

After running the worker front, you can access the worker front with visiting `http://localhost:8222`.
To run the function `W1`, request `http://localhost:8222/execute?name=W1`.
