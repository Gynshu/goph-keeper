# Goph-Keeper
This is a password manager written in Go.
Graduation project of Advanced golang course at Yandex.Praktikum
##  Tech stack
- [Go](https://golang.org/)
- [MongoDB](https://www.mongodb.com/)
- [Docker-compose](https://docs.docker.com/compose/)
- [TOTP](https://en.wikipedia.org/wiki/Time-based_One-time_Password_algorithm)
- [AES](https://en.wikipedia.org/wiki/Advanced_Encryption_Standard)
- [Chi router](https://github.com/go-chi/chi/v5)
- [Resty](https://github.com/go-resty/resty/v2)
- [Zero log](https://github.com/rs/zerolog/log)
## Installation

```bash
git clone https://github.com/gynshu-one/goph-keeper
cd goph-keeper
```
### Server part
Inside the project directory, to build and run server:
```bash
make serve
```
Under the hood is this:
```bash
@echo "Starting server..."
sudo chmod +x server/cert/gen_keys.sh
cd server/cert && ./gen_keys.sh
# runs mongodb
docker-compose up --build -d
# builds the server
go build -o ./server_cmd ./server/cmd/main.go
# runs the server
./server_cmd -c ./server/config.json
```

### Client part
Inside the project directory, to build and run client:
```bash
make ui
```

Under the hood is this:
```bash
@echo "Building client..."
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd -c "./client/config.json"
```


see [Makefile](https://github.com/gynshu-one/goph-keeper/Makefile)
This will generate TSL certs build docker image and run it on port 8080 as well as mongodb on port 27017


## Signing up and logging in
<img style="max-width:600px" src="https://i.imgur.com/DDsrsM6.png">

For simplicity, the password and secret are stored in the OS keychain.
Every other log in would grab username from  /temp/session_id file and check OS keychain for password and secret

After logging in to a server with only password (Not secret) you will receive session cookie for 24 hours.

## Configuring

The client will read the config.json file from its working directory.
```json
{
  "SERVER_IP": "localhost:8080"
}
```

The Server will config.json from its working dir
```json
{
  "MONGO_URI": "mongodb://admin:password@mongo_db:27017",
  "HTTP_SERVER_PORT": "8080",
  "CERT_FILE_PATH": "common/cert/cert.pem",
  "KEY_FILE_PATH": "common/cert/key.pem"
}
```

### Basic ui

#### Main page

<img style="max-width:600px" src="https://i.imgur.com/EswW6Xo.png">

Adding a new bank card

<img style="max-width:600px" src="https://i.imgur.com/hXx4UzS.png">

#### Editing it

<img style="max-width:600px" src="https://i.imgur.com/AqI3rRM.png">

#### New login

<img style="max-width:600px" src="https://i.imgur.com/LdGKuON.png">

#### Deletion

<img style="max-width:600px" src="https://i.imgur.com/6OIPMR7.png">

#### Generating One time

<img style="max-width:600px" src="https://i.imgur.com/UC2Fi5W.png">


<img style="max-width:600px" src="https://i.imgur.com/bMAYiCI.png">

## API

Server has 4 endpoints
Which are defined in [router.go](https://github.com/gynshu-one/goph-keeper/server/api/router/router.go)
### /user/create
Creates new user with username and password with url params
```
https://localhost:8080/user/create?email=your_username&password=your_password
```
### /user/login
Logs in user
creates session cookie for 24 hours
```
https://localhost:8080/user/login?email=your_username&password=your_password
```
### /user/logout
Logs out user and deletes session cookie
```
https://localhost:8080/user/logout
```
### /user/sync
Synchronizes user data with server. Server check if user has session cookie via
[middleware.go](https://github.com/gynshu-one/goph-keeper/server/api/middlewares/middleware.go)
which uses [session.go](https://github.com/gynshu-one/goph-keeper/server/api/session/session.go)
```
https://localhost:8080/user/sync
```
`POST` data slice of `DataWrapper` structs
```go
// DataWrapper is a struct that wraps BasicData and provides additional information about the data
// such as owner id, type, name, updated_at, created_at, deleted_at
// it makes easier to store data in the database that shouldn't know anything about the data
type DataWrapper struct {
// ID is the unique identifier of the data
ID string `json:"id" bson:"_id"`
// OwnerID is the user who owns this data
OwnerID string `json:"owner_id" bson:"owner_id"`
// Type is the type of the data such as ArbitraryTextType, BankCardType, BinaryType, LoginType
Type      string `json:"type" bson:"type"`
Name      string `json:"name" bson:"name"`
UpdatedAt int64  `json:"updated_at" bson:"updated_at"`
CreatedAt int64  `json:"created_at" bson:"created_at"`
DeletedAt int64  `json:"deleted_at" bson:"deleted_at"`
// This is the actual data that is stored in the database
// Encrypted with user's secret
Data      []byte `json:"data" bson:"data"`
}
```

struct in [general.go](https://github.com/gynshu-one/goph-keeper/common/models/general.go)

`Sync` happens imminently after logging in and every item creation, deletion or update.

`Client` does not have cache of the data, it only stores session cookie and username in /temp/session_id file.


Deleted items would force server to remove `DataWrapper`'s Data field and set `DeletedAt` field.