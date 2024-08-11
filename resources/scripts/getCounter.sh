#!/bin/bash

for ((i=1; i<=10000; i++))
do
  curl -X GET 'localhost:8080/counter' &
done
