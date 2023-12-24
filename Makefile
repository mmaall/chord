http-server:
	go build -o http-server -a cmd/http_server/main.go

http-server-al2:
	GOARCH=amd64 GOOS=linux go build -o http-server -a cmd/http_server/main.go

http-server-clean:
	rm -f http-server

images: 
	docker build --tag http-server ./ -f dockerfiles/HttpServer.docker

images-clean:
	docker image rm http-server -f 

test:
	go test -v ./...

clean: http-server-clean images-clean

