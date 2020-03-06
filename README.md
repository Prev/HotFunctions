# HotFunctions [![Build Status](https://travis-ci.org/Prev/HotFunctions.svg)](https://travis-ci.org/Prev/HotFunctions) [![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Prev/HotFunctions/blob/master/LICENSE)


A simple FaaS framework and tools to test performance of load balancing algorithms.
Consists of 5 parts like described in below:

- [Worker front](./worker_front): Codes for the worker node to run the user functions in the Docker container.
- [Load balancer](./load_balancer): Codes for distributing requests to multiple worker nodes
- [Simulator](./simulator): Codes to run a simulation with load balancer
- [Sample functions](./sample_functions): Pre-defined user functions to run a simulation
- [Utils](./utils): Utilities like generating simulation scenario or analyzing the result.


## Run the project
### Prerequisites

- [Golang](https://golang.org/) >= 1.13 (Since we use `gomod` for dependency control, at least 1.11 is required)
- [Docker](https://www.docker.com/) (Functions are running in the Docker container)


### Step1: Configure a worker node (or worker nodes)

Follow the "How to run" section on [worker front page](./worker_front).  
After configuration, visit http://localhost:8222/execute?name=W1 to test worker node.

Our framework supports 30 different functions.
You can see the detail of the supporting sample functions in [sample functions page](./sample_functions).

### Step2: Configure a load balancer

Follow the "How to run" section on [load balancer page](./load_balancer).  
After configuration, visit http://localhost:8111/execute?name=W1 to test load balancer.
Note that you should not close the worker node while running the load balancer.


### Step3: Run simulator and analyze result with utils

Follow the "How to run" section on [simulator page](./simulator).
After the simulation is finished, utils for analyzing the logs are prepared on [utils](./utils).
