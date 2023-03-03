BINARY_NAME=kview.app

## build: build binary and package app
build:
	rm -rf ${BINARY_NAME}
	~/go/bin/fyne package -os darwin -release

## run: builds and runs the application
run:
	go run .

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${BINARY_NAME}
	@echo "Cleaned!"

## test: runs all tests
test:
	go test . -v