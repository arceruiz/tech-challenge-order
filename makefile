run-db:
	docker-compose -f deployments/db-docker-compose.yml up -d

run-tests:
	go test $$(go list ./... | grep -v /data/) -coverprofile=cover.out.tmp && cat ./cover.out.tmp | grep -v "mock.go" > ./cover.out && go tool cover -html=cover.out 

test:
	go test $$(go list ./... | grep -v /data/) -coverprofile=cover.out.tmp

test-build-bake:
	docker build -t mauricio1998/order-service . -f build/Dockerfile

run-localstack:
	docker run --rm -it -p 4566:4566 localstack/localstack

run-infra: connect-localstack create-order-queue create-payment-queue create-payment-payed-queue

connect-localstack:
	awslocal kinesis list-streams

create-order-queue:
	awslocal sqs create-queue --queue-name orderqueue

create-payment-queue:
	awslocal sqs create-queue --queue-name paymentpendingqueue

create-payment-payed-queue:
	awslocal sqs create-queue --queue-name paymentpayedqueue

create-payment-cancelled-queue:
	awslocal sqs create-queue --queue-name paymentcancelledqueue

test-build-bake:
	docker build -t docker.io/mauricio1998/order-service . -f build/Dockerfile

docker-push:
	docker push docker.io/mauricio1998/order-service

boiler-plate: test-build-bake docker-push
