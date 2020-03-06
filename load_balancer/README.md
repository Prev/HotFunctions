# Load Balancer

Distribute functions to worker nodes. There are 5 different scheduling algorithms in the framework like in the below:

- Round Robin
- Least Loaded: Equivalent with least connected of NGINX
- Consistent Hashing with bounded load: Variation of Consistent Hashing algorithm that limits the load per each node [arxiv](https://arxiv.org/abs/1608.01350)
- PASch extended: An extended version changing the package id to an application name of the [PASch algorithm](https://ieeexplore.ieee.org/document/8752939)
- Ours: Proposing algorithm that maximizes locality while minimizing load imbalance

### Configuring worker nodes

We recommended you to configure two or more worker nodes to use the load balancer correctly.
After [configuring worker nodes](../worker_front), edit `nodes.config.json` before starting the load balancer.
Note that the http protocol and port number should be included in the file.

### How to run

```bash
go run *.go rr|ll|ch|pasch|ours
```

Options `rr`, `ll`, `ch`, `pasch`, `ours` mean Round Robin, Least Loaded, Consistent Hashing with bounded load, PASch extended, and Ours, respectively.


Load balancer will send a similar response with worker node, but additional information from load balancing will be appended to the response like below.

```json
{
	"Result": {
		"statusCode": 200,
		"body": "Hello World2999663"
	},
	"ExecutionTime": 929,
	"InternalExecutionTime": 98,
	"Meta": {
		"ImageBuilt": false,
		"UsingPooledContainer": false,
		"UsingExistingRestContainer": false,
		"ContainerName": "hf_w1__9572_27887",
		"ImageName": "hf_w1"
	},
	"LoadBalancingInfo": {
		"WorkerNodeId": 0,
		"WorkerNodeUrl": "http://localhost:8222",
		"Algorithm": "ours",
		"AlgorithmLatency": 0.015869140625
	}
}
```