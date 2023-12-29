# Demos - Testing and Proving That Temporal Does What It Says
Playing with new frameworks is fun, and Temporal is a joy in particular because it can focus you to work on the core concepts and code, and ignore things not on the happy path.

## Simple Demo
This is the basic demo described in the [setup instructions](./setup.md). 
You can see an order get processed, maybe fail randomly.
1. Start the worker :
```shell
go run workers/main.go
```

2. Test with :
```shell 
go run starter/main.go
```


## Killing the Process Doesn't Make Anything Break
This demo shows that Temporal applications can survive process crashes. It will work with any order number, but a pause is built in for [#37005](./workflows/process_order.go)
See the [kill_worker demo script](./demoscripts/kill_worker.sh).

1. Start worker in a terminal:
```shell 
go run workers/main.go
```
2. New terminal:
```shell 
chmod +x ./demoscripts/*
cd ./demoscripts/
./kill_worker.sh
```
3. Wait until the order gets submitted, and then kill the worker
4. Observe that the starter is happily waiting, and the Temporal UI shows the workflow still alive
5. Start the worker again, and it will pick up where it left off and complete
6. Tada, your code is pretty bulletproof: the order completed, no steps were duplicated, and you didn't have to do anything besides use Temporal to do it.

## Different Order Number As Input
If you want to avoid duplicate order number checks, you can pass in an order number:
1. Start worker in a terminal if you haven't already:
```shell 
go run workers/main.go
```
2. New terminal, changing the first argument to your favorite order number:
(Don't change the item number, we only have the one item - but feel free to change the quantity for fun!)
```shell 
go run starter/main.go 11235813 123456 1 VISA-123456
```

## Errors, Errors Everywhere
1. Modify [is_error.go](./utils/is_error.go) to have a higher likelihood of error:
```go
	// throw errors 1 time out of 100
	if rand.Intn(100) > 99 {
		return true
	}
```
to 
```go
	// throw errors 34 time out of 100
	if rand.Intn(100) > 66 {
		return true
	}
```

2. Start worker in a terminal:
```shell 
go run workers/main.go
```
3. New terminal:
```shell 
go run starter/main.go
```

## Duplicate Order Idempotence
This shows how to implement idempotence in Temporal. Check out Pierre's excellent [original code](https://github.com/PierreSylvain/idempotence) and full blog post here: [Idempotence in Temporal.io, a Look into Technical Architectures](https://medium.com/@ps.augereau/idempotence-in-temporal-io-a-look-into-technical-architectures-11d20a0fc860)

1. Start worker in a terminal:
```shell 
go run workers/main.go
```
2. New terminal:
```shell 
chmod +x ./demoscripts/*
cd ./demoscripts/
./idempotencydemo.sh
```
3. Observe that two orders with the same ID are submitted.
4. Observe that the idempotency of the process is handled: the order is processed once, no matter the errors it incurs during processing.
5. Observe that the inventory did not change in the second order.
