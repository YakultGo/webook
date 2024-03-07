.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f yakultd/webook:v0.0.1 || true
	@docker build -t yakultd/webook:v0.0.1 .