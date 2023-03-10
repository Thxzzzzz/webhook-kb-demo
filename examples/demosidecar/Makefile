IMG ?= 379809513/demo-sidecar:latest

run:
	go run main.go

docker-build:
	docker build -t ${IMG} .

docker-push:
	docker push ${IMG}

docker-run:
	docker run --name demo-sidecar ${IMG} --port 9000:9000