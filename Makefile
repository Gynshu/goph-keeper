serve:
	@echo "Starting server..."
	sudo chmod +x common/cert/gen_keys.sh
	cd common/cert && ./gen_keys.sh
	docker-compose up --build -d

ui:
	@echo "Building client..."
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd