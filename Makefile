run-server:
	go run cmd/server/main.go

get-counter:
	sh ./resources/scripts/getCounter.sh

health-check:
	sh ./resources/scripts/healthCheck.sh
