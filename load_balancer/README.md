## Load Balancer

Distribute functions to worker nodes. There are 4 different scheduling algorithms in the framework like in the below:

- Least Loaded: Scheduler picks the node who has minimum executing tasks
- Consistent Hashing: Scheduler picks the node by Consistent Hashing algorithm where key is the function name
- Ours_static: Static scheduler with Lookup-table without reassignment
- Ours_adaptive: Adaptive scheduler with lookup-table and reassignment

### How to run

```bash
go run *.go ll|ch|tb|ad
```


```bash
go get -u github.com/lafikl/consistent
```
