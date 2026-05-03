ifeq ($(OS),Windows_NT)

build:
	go build -o bin\cpdi.exe cpdi\cmd\cpdi

test:
	go test cpdi\cmd\cpdi

else

build:
	go build -o ./bin/cpdi cpdi/cmd/cpdi

test:
	go test -v cpdi/cmd/cpdi

endif