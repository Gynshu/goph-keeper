start_server:
	@echo "Starting server..."
	sudo chmod +x shared/cert/gen_keys.sh
	cd shared/cert && ./gen_keys.sh
	@#go run main.go

run:
	@echo "Building client..."
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd