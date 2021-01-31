# Simulator

Start a simulation with load balancer and event stream.


### How to run

Run like `go run *.go <url> [<event_file>] [<output_file>]`.
Default values for `event_file` and `output_file` are `data/events.csv` and `logs/YYYY-MM-DD/HH:ii:ss.log` respectively.

Notes that default simulation file(`data/events.csv`) has a lot of events, so you may need multiple worker nodes.
We recommand at least 8 worker nodes with 4 vCPU for each node (e.g. EC2 m4.xlarge).

```bash
$ go run *.go http://localhost:8111
$ go run *.go http://localhost:8111 data/events_small.csv
$ go run *.go http://localhost:8111 data/events_small.csv output.log
```

If you want to generate your own event stream, you can use `generate_event_streams.py` on [utils](../utils).
The CSV file only contains `FunctionName, StartTime` pairs, so you can easily create another simple simulation scenario.
