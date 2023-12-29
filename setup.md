# Setup
You need to have a Temporal server up and running and Go installed. This currently uses temporal cloud but you can configure it to use local.

CLone the repo :
```shell
git clone https://github.com/joshmsmith/temporal-idempotence-by-validation
```

## Temporal Cloud configuration
This example assumes that you have a temporal cloud configured and have local client certificate files for your namespace.
The values are passed into the demo app using environment variables, example direnv .envrc file is included in the repo:

```
# direnv .envrc

# Temporal Cloud connection
# region: us-east-1
export TEMPORAL_HOST_URL="myns.abcdf.tmprl.cloud:7233"
export TEMPORAL_NAMESPACE="myns.abcdf"
#export TEMPORAL_HOST_URL="localhost:7233"
#export TEMPORAL_NAMESPACE="workflow-web"

# If self-hosted, skip TLS certs
export USE_TLS=true
#export USE_TLS=false

# tclient-myns client cert
export TEMPORAL_TLS_CERT="/Users/myuser/.temporal/tclient-myns.pem"
export TEMPORAL_TLS_KEY="/Users/myuser/.temporal/tclient-myns.key"

# Optional: path to root server CA cert
export TEMPORAL_SERVER_ROOT_CA_CERT=
# Optional: Server name to use for verifying the server's certificate
export TEMPORAL_SERVER_NAME=

export TEMPORAL_INSECURE_SKIP_VERIFY=false


# App temporal taskqueue names
export TICKET_ORDER_MANAGEMENT_TASK_QUEUE="idempotence-by-validation"

# timer for transfer table to be checked (seconds)
export CHECK_TRANSFER_TASKQUEUE_TIMER=20
# timer for demo delay between Withdraw and Deposit Activities
export DELAY_TIMER_BETWEEN_WITHDRAW_DEPOSIT=15

# payload data encryption
export ENCRYPT_PAYLOAD=false
export DATACONVERTER_ENCRYPTION_KEY_ID=mysecretkey

# Set to enable debug logger logging
export LOG_LEVEL=debug

# local JSON backend db connection
export DATABASEPATH=database/




```


Then start the worker :
```shell
go run workers/main.go
```

And finally, test with :
```shell 
go run starter/main.go
```

## Other Things To Try
Check out more stuff in [demos.md](./demos.md)

