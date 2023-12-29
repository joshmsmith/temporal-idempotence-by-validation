# Temporal Idempotence By Validation Demo
This is a simple [Temporal](https://temporal.io/) demo in Go that demonstrates **idempotent by validation**. Sometimes with Temporal, you need to call a service that isn't [idempotent](https://en.wikipedia.org/wiki/Idempotence): calling it more than once with the same inputs is bad - may cause duplicate entries or other unintended behavior. 
In such a case, you can set a policy to call the non-idempotent activity only once, and then **validate** that it was successful, and retry it if not. In this way, the *process* of creation-validation becomes idempotent; it won't run the non-idempotent activity again unless it is needed, and the process can be called as many times as needed until it succeeds.

This demo demonstrates this *idempotency by validation* in the context of a ticket ordering system.

Here is a sample ticket order process:
**Input:** order ID
1. Get token for calling the ticket reservation system
2. Retrieve reservation by order
3. Retrieve payment info 
4. Create ticket for reservation and payment info

**Output:** one and only one ticket created in database, return ticket number

Step #4, *create ticket*, is *not* idempotent. It should only be called one time. However, we can create a step that *validates that a ticket exists*. Let's make a process like this: 

1. Get token for calling the ticket reservation system
2. Retrieve reservation by order
3. Durably create tickets:
    a. Retrieve payment info 
    b. Create ticket for reservation and payment info -- may fail and can only be called once
    c. Validate ticket for reservation
    d. Loop here until a ticket is created

Step #3 is our *idempotency loop*. It could be a [child workflow](https://docs.temporal.io/workflows#child-workflow) and leverage Temporal's retryability, but I kept it as a loop for readability's sake.

This demo additionally demonstrates the [durability](https://temporal.io/how-it-works) of a process implemented in Temporal:
1. Crashing the process doesn't kill it. Upon resume it picks up right where it left off.
2. Errors are recovered without thought or work
3. At-Least-Once execution: activities succeed at least once per the workflow

These capabilities are **great to develop with** and **change the way I think about doing development.** 

![durable_execution](./resources/durable_execution_abstraction_small.png)

As a developer I can **focus** just on what I want to do, and Temporal manages what happens when things don't work out. 

While working on this project, I created **many** bugs in my activities, and all I had to do to fix my in-flight orders was fix the code bugs and restart the worker process. The  errors went away and none of the workflows failed, they all succeeded.

*Zero workflow processes failed in the building of this demo*.


## Process Results
The code in [starter](./starter/main.go) demonstrates the workflow. Initially, there are only a couple sample tickets in our "[database](./database/)". After ticket creation, you will see more created there.

Here is a sample ticket:

```json
{
    "orderID": "order-112358",
    "ticketID": "TICKET-42618",
    "paymentInfo": "VISA-5197-988-3381-2526"
}
```

The credit card number is random and not a valid credit card number. In addition, because the card info is handled all in one activity, it never reaches Temporal's servers, which you can verify by looking at the workflow JSON.


# Getting Started
See [Setup Instructions](./setup.md).

### First Demo
After the setup is done, you can do the  basic demo described in the [setup instructions](./setup.md). 
You can see an order get processed, maybe fail randomly.
1. Start the worker :
```shell
go run workers/main.go
```

2. Test with :
```shell 
go run starter/main.go
```

# Next Steps
1. Check out the ways to [demonstrate that this works nicely](./demos.md)
2. Play around with the code in new ways, try to break Temporal, maybe try some [retry policies](https://docs.temporal.io/retry-policies#:~:text=A%20Retry%20Policy%20works%20in,or%20an%20Activity%20Task%20Execution.) 
3. Feel free to fork and contribute!