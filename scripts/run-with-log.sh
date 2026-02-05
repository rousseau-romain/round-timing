#!/bin/sh
mkdir -p logs
exec go run -buildvcs . 2>&1 | tee -a ./logs/app.log
