serve:
	@echo "Starting server..."
	# generates self-signed certificate
	mkdir -p ~/.goph-keeper
	openssl req -x509 -newkey rsa:4096 -keyout ~/.goph-keeper/key.pem -out ~/.goph-keeper/cert.pem -days 365 -nodes -subj '/CN=localhost' -addext 'subjectAltName = DNS:localhost'
	# runs mongodb
	docker-compose up --build -d
	# builds the server
	go build -o ./server_cmd ./server/cmd/main.go
	# runs the server with the pre-generated certificate and predefined mongo URI
	./server_cmd -port 8080 -cert ~/.goph-keeper/cert.pem -key ~/.goph-keeper/key.pem -mongo_uri mongodb://admin:password@localhost:27017

ui:
	@echo "Building client..."
	mkdir -p ~/.goph-keeper
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd -addr localhost:8080 -poll 5s -dump 10s

tests:
	@echo "Running tests..."
#	sudo chmod +x server/cert/gen_keys.sh
#	cd server/cert && ./gen_keys.sh
	docker-compose down
	docker-compose up --build mongo_db -d
	go test -v ./...