all: deps coverall

deps:
	go get -t ./...
	go get golang.org/x/tools/cmd/cover
	go get github.com/mattn/goveralls
	#go get github.com/golang/lint/golint

test:
	go test -race -cover ./...

coverall:
    go test -v -covermode=count -coverprofile=coverage.out \
    $HOME/gopath/bin/goveralls -coverprofile=coverage.out \
    -service=travis-ci -repotoken ODVVqrmKEbXYriRSKushcupIC0UUwjBxv

validate: lint
	go vet ./...
	test -z "$(gofmt -s -l . | tee /dev/stderr)"

lint:
	out="$$(golint ./...)"; \
	if [ -n "$$(golint ./...)" ]; then \
		echo "$$out"; \
		exit 1; \
	fi
