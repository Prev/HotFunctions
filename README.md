# HotFunctions

A simple FaaS framework and tools to test the performance of load balancing algorithms.
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
Visit http://localhost:8222/execute?name=W1 to test worker node (This command will execute function `W1` and send you a response).
Reponse would be like below:

```json
{
	"Result": {
		"statusCode": 200,
		"body": "Hello World2998950"
	},
	"ExecutionTime": 3373,
	"InternalExecutionTime": 98,
	"Meta": {
		"ImageBuilt": true,
		"UsingPooledContainer": false,
		"UsingExistingRestContainer": false,
		"ContainerName": "hf_w1__9437_98081",
		"ImageName": "hf_w1"
	}
}
```

Our framework supports 30 different functions.
You can see the detail of the supporting sample functions in [sample functions page](./sample_functions).

### Step2: Configure a load balancer

Follow the "How to run" section on [load balancer page](./load_balancer).
Visit http://localhost:8111/execute?name=W1 to test load balancer.
Load balancer will send a similar response with worker node, but additional information from load balancing will be appended to the response like below.

```json
{
	"Result": {
		"statusCode": 200,
		"body": "Hello World2999663"
	},
	"ExecutionTime": 929,
	"InternalExecutionTime": 98,
	"Meta":{
		"ImageBuilt": false,
		"UsingPooledContainer": false,
		"UsingExistingRestContainer": false,
		"ContainerName": "hf_w1__9572_27887",
		"ImageName": "hf_w1"
	},
	"LoadBalancingInfo":{
		"WorkerNodeId": 0,
		"WorkerNodeUrl": "http://localhost:8222",
		"Algorithm": "ours",
		"AlgorithmLatency": 0.015869140625
	}
}
```

### Step3: Run simulator and analyze result with utils

Follow the "How to run" section on [simulator page](./simulator).
After the simulation is finished, utils for analyzing the logs are prepared on [utils](./utils).
