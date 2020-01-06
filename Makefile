GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=cryp

build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/cryp

db:
	docker build -t app . --rm=true
	docker run -d -p 8000:8000 app
	docker rmi $(docker images -f "dangling=true" -q)

dp:
	docker stop $(docker container ps -a -q)
	docker rm -f $(docker ps -a -q)
	docker image prune -a -f

run:
	$(GOCMD) ./$(BINARY_NAME)

vue:
	cd ./ui
	npm run build
