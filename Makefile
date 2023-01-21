run: 
	# Run app and use the root directory to open the repository
	go run cmd/gitglance/main.go .

test:
	 go test ./...