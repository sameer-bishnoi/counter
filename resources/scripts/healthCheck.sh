#!/bin/bash

for ((i=1; i<=10; i++))
do
  curl -X GET 'localhost:8080/health-check' &
done
