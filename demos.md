# Demos - Testing and Proving That Temporal Does What It Says
Playing with new frameworks is fun, and Temporal is a joy in particular because it can focus you to work on the core concepts and code, and ignore things not on the happy path.

## Simple Demo
This is the basic demo described in the [setup instructions](./setup.md). 
You can see an ticket order get processed, and fail randomly.
1. Start the worker :
```shell
go run workers/main.go
```

2. Test with :
```shell 
go run starter/main.go
```


## Killing the Process Doesn't Make Anything Break
This demo shows that Temporal applications can survive process crashes. 

1. Start worker in a terminal:
```shell 
go run workers/main.go
```
2. New terminal:
```shell 
go run starter/main.go
```
3. Wait until the workflow is going, and then kill the worker
4. Observe that the starter is happily waiting, and the Temporal UI shows the workflow still alive
5. Start the worker again, and it will pick up where it left off and complete
6. Tada, your code is pretty bulletproof: the order completed, no steps were duplicated, and you didn't have to do anything besides use Temporal to do it.

## Errors, Errors Everywhere
1. Modify the various functions in the [ticket system](./ticket/ticket_system.go) with different error options from [is_error.go](./utils/is_error.go) to have a higher likelihood of error:
```go
	func IsError() bool {
	// throw errors 1 time out of 100
	if rand.Intn(100) < 99 {
		return false
	}
	return true
}

func IsErrorMoreLikely() bool {
	// throw errors 10 times out of 100
	if rand.Intn(100) < 90 {
		return false
	}
	return true
}

func IsErrorPrettyLikely() bool {
	// throw errors 50 times out of 100
	if rand.Intn(100) < 50 {
		return false
	}
	return true
}


func IsErrorVeryLikely() bool {
	// throw errors 80 times out of 100
	if rand.Intn(100) < 20 {
		return false
	}
	return true
}
```
see for example:
```go
	// simulate a random error before creating a ticket
	if utils.IsErrorPrettyLikely() {
		return "", errors.New("CREATE-TICKET-ERROR")
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

## Idempotence
1. Make sure the errors in CreateTicket() are very likely, then run the demo:
2. Start worker in a terminal:
```shell 
go run workers/main.go
```
3. New terminal:
```shell 
go run starter/main.go
```
Observe that despite high frequency of errors in CreateTicket, it only ever makes one ticket per order/workflow.
Even if everything else has a high frequency of errors, it will still run once and only once.

