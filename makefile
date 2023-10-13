build:
	@go build -o bin/gopaste.exe

run: build
	@bin/gopaste.exe