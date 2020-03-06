# HotFunctions [![Build Status](https://travis-ci.org/Prev/HotFunctions.svg)](https://travis-ci.org/Prev/HotFunctions) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Prev/HotFunctions/blob/master/LICENSE)


A simple FaaS framework and tools to test performance of load balancing algorithms.
Consists of 5 parts like described in below:

- [Worker front](./worker_front): Codes for the worker node to run the user functions in the Docker container.
- [Load balancer](./load_balancer): Codes for distributing requests to multiple worker nodes
- [Simulator](./simulator): Codes to run a simulation with load balancer
- [Sample functions](./sample_functions): Pre-defined user functions to run a simulation
- [Utils](./utils): Utilities like generating simulation scenario or analyzing the result.


## Getting Started
### Prerequisites

- [Golang](https://golang.org/) >= 1.13 (Since we use `gomod` for dependency control, at least 1.11 is required)
- [Docker](https://www.docker.com/) (Functions are running in the Docker container)


### Step1: Configure a worker node (or worker nodes)

```bash
$ cd worker_front
$ go run *.go start
```

You can see the details in [worker front](./worker_front) page.  
After configuration, visit http://localhost:8222/execute?name=W1 to test worker node.

Our framework supports 30 different functions.
You can see the detail of the supporting sample functions in [sample functions page](./sample_functions).

### Step2: Configure a load balancer

```bash
$ cd load_balancer
$ go run *.go rr|ll|ch|pasch|ours
```

You can also see the details on [load balancer](./load_balancer) page.  
After configuration, visit http://localhost:8111/execute?name=W1 to test load balancer.
Note that you should not close the worker node while running the load balancer.


### Step3: Run simulator and analyze result with utils

Follow the "How to run" section on [simulator page](./simulator).
After the simulation is finished, utils for analyzing the logs are prepared on [utils](./utils).
