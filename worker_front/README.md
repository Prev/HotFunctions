# Worker Front

Task runner of worker node.
Receives application request with RESTful API and executes the function with Docker.

### Base features

- Build image automatically if image not exists
	- Download the target function source code (from S3)
	- Attach environment by function's demands (python3.7, NodeJS 12, Java 8)
	- Build docker image with the function source code
- Creates containers for function execution
	- After execution, remove container and send responses

### Optional features
1. Limit the cache size of downloaded user codes
	- In default mode, there is no limitation of built images and cached codes
	- With this option, worker node manages Docker image and user code with LRU policy
	- To use this option, set environment variable `USER_CODE_SIZE_LIMIT` with proper size.
		- For example, to set limit as 100MB, you can set this variable as `100000000`

2. Provides pre-warmed container pool
	- Based on the idea to reduce the preparation time by pre-building the environment before the function is executed
	- To use this option, you need to configure environment variables `CONTAINER_POOL_LIMIT` and `CONTAINER_POOL_NUM`.
		- Worker node will generate pre-warmed containers by LRU policy
		- `CONTAINER_POOL_LIMIT` means the number of caching applications
		- `CONTAINER_POOL_NUM` means the number of pre-warmed containers per application (default: 6)

3. Supports two-levels of isolation
	- Based on the idea to let share the sandbox among same applications
	- To use this option, set `USING_REST_MODE` environment variable as `true`
	- The default lifetime of each idle container is 10 seconds. You can change is value with `REST_CONTAINER_LIFE_TIME`.


### How to run

We support two arguments on starting worker front.
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


### Configure with RESTful API

You can also change configuration after starting the worker front.
Visit `http://localhost:8222/configure` and change options by passing queries like `?using_rest_mode=true&rest_container_life_time=20`.

