serve:
	@echo "Starting server..."
	sudo chmod +x server/cert/gen_keys.sh
	cd server/cert && ./gen_keys.sh
	# runs mongodb
	docker-compose up --build -d
	# builds the server
	go build -o ./server_cmd ./server/cmd/main.go
	# runs the server
	./server_cmd -c ./server/config.json

ui:
	@echo "Building client..."
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd -c "./client/config.json"

tests:
	@echo "Running tests..."
#	sudo chmod +x server/cert/gen_keys.sh
#	cd server/cert && ./gen_keys.sh
	docker-compose down
	docker-compose up --build mongo_db -d
	go test -v ./...