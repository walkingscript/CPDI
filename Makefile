build:
	go build -o ./bin/cpdi cpdi/cmd/cpdi

test:
	go test cpdi/cmd/cpdi
