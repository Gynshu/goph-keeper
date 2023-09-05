# Goph-Keeper
This is a password manager written in Go.
Graduation project of Advanced golang course at Yandex.Praktikum

## Installation

inside the project directory run:
```bash
make serve
```
under the hood is this:
```bash
@echo "Starting server..."
sudo chmod +x common/cert/gen_keys.sh
cd common/cert && ./gen_keys.sh
docker-compose up --build -d
```
see [Makefile](https://github.com/gynshu-one/goph-keeper/Makefile)
This will generate TSL certs build docker image and run it on port 8080 as well as mongodb on port 27017

## Getting started with

You'll need to have Go installed on your machine. You can download it [here](https://golang.org/dl/).
To run client

```bash
make ui
```
under the hood is this:
```bash
@echo "Building client..."
go build -o ./client_cmd ./client/cmd/main.go
./client_cmd
```
## Signing up and logging in
<img style="max-width:400px" src="https://imgur.com/DDsrsM6">
For simplicity, the password is and secret is store in OS keychain
Every other log in would grab username from  /temp/session_id file and check OS keychain for password and secret

After logging in to a server with only password (Not secret) you will receive session cookie for 24 hours.

## Configuring

Client will read config.json file from its base dir which looks like this
```json
{
  "SERVER_IP": "localhost:8080",
  "POLL_TIMER": 5,
  "DUMP_TIMER": 10,
  "CACHE_FOLDER": "/data-keeper"
}
```

### Basic ui

Main page
<img style="max-width:400px" src="https://imgur.com/EswW6Xo">

Adding new bank card
<img style="max-width:400px" src="https://imgur.com/hXx4UzS">

Editing it
<img style="max-width:400px" src="https://imgur.com/AqI3rRM">

New login
<img style="max-width:400px" src="https://imgur.com/LdGKuON">

Deletion
<img style="max-width:400px" src="https://imgur.com/6OIPMR7">

Generating One time
<img style="max-width:400px" src="https://imgur.com/UC2Fi5W">
<img style="max-width:400px" src="https://imgur.com/bMAYiCI">

## API

Server has 4 endpoints
Which are defined in [router.go](https://github.com/gynshu-one/goph-keeper/server/api/router/router.go)
### /user/create
Creates new user with username and password with url params
```
https://localhost:8080/user/create?username=your_username&password=your_password
```
### /user/login
Logs in user
creates session cookie for 24 hours
```
https://localhost:8080/user/login?username=your_username&password=your_password
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
with post data list defined in `DataWrapper`
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

Every time user creates something new or signs in it sends all data to server and server updates it in database if needed.
Deleted items would force server to remove `DataWrapper`'s Data field and set `DeletedAt` field.