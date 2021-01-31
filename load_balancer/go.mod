module github.com/Prev/HotFunctions/load_balancer

go 1.13

require (
	github.com/Prev/HotFunctions/load_balancer/scheduler v1.0.0
	github.com/Prev/HotFunctions/worker_front/types v1.0.0
)

replace (
	github.com/Prev/HotFunctions/load_balancer/scheduler => ./scheduler
	github.com/Prev/HotFunctions/worker_front/types => ../worker_front/types
)
