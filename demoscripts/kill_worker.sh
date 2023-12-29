#!/bin/bash
echo "Submitting order 37005 which will cause a delay in processing."
echo "Prepare to kill the worker, then view in the Temporal UI, then start the worker again."
cd ..
go run starter/main.go 37005 123456 999 VISA-5551212
echo "It completed, nice"
cd -