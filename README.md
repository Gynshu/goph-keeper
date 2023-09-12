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
- [UI with Tveiw](https://github.com/rivo/tview)
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
# generates self-signed certificate
mkdir -p ~/.goph-keeper
openssl req -x509 -newkey rsa:4096 -keyout ~/.goph-keeper/key.pem -out ~/.goph-keeper/cert.pem -days 365 -nodes -subj '/CN=localhost' -addext 'subjectAltName = DNS:localhost'
# runs mongodb
docker-compose up --build -d
# builds the server
go build -o ./server_cmd ./server/cmd/main.go
# runs the server with the pre-generated certificate and predefined mongo URI
./server_cmd -port 8080 -cert ~/.goph-keeper/cert.pem -key ~/.goph-keeper/key.pem -mongo_uri mongodb://admin:password@localhost:27017
```
`port` is the server port default: 8080<br>
`cert` is the path to the certificate default: cert/cert.pem<br>
`key` is the path to the key default: cert/key.pem<br>
`mongo_uri` is the mongodb uri default: mongodb://admin:password@localhost:27017<br>
if you run server without any flag or without cert and key it will generate self-signed certificate for `localhost` and run on port 8080


### Client part
Inside the project directory, to build and run client:
```bash
make ui
```
Under the hood is this:
```bash
@echo "Building client..."
	go build -o ./client_cmd ./client/cmd/main.go
	./client_cmd -addr localhost:8080 -poll 5s -dump 10s
```
`addr` is the server ip address default: localhost:8080<br>
`poll` is the timer for synchronizing data with server default: 5s<br>
`dump` is the timer for dumping data to the server default: 10s<br>

You can also run client without any flags and it will use default values

see [Makefile](https://github.com/gynshu-one/goph-keeper/blob/main/Makefile)

## Docs
```bash
 go install golang.org/x/tools/cmd/godoc
 godoc -http=:6060
```
then visit:<br>
http://localhost:6060/pkg/github.com/gynshu-one/goph-keeper/


## Signing up and logging in
<img style="max-width:600px" src="https://i.imgur.com/DDsrsM6.png">

For simplicity, the password and secret are stored in the OS keychain.
Every other log in would grab first user from OS keychain for password and secret

After logging in to a server with only password (Not secret) you will receive session cookie for 24 hours.

## Configuring

The client will read the config.json file from its working directory.
```json
{
  "SERVER_IP": "localhost:8080"
}
```

The Server will read config.json from its working dir
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
This page is basically a list of items, each item takes two rows first is a name of an item second is a type of it.

If item was deleted recently client would see name and "deleted" message. 

Header part contains item creation buttons

If you click on item name you will be switched to item edition page, it looks just like creation page, except additional buttons "delete" or item's `type` specific "generate one time" for `login` type
<img style="max-width:600px" src="https://i.imgur.com/EswW6Xo.png">

### Example of new item creation page

<img style="max-width:600px" src="https://i.imgur.com/hXx4UzS.png">

#### Edit page

<img style="max-width:600px" src="https://i.imgur.com/AqI3rRM.png">

#### New login

<img style="max-width:600px" src="https://i.imgur.com/LdGKuON.png">

#### Deletion

<img style="max-width:600px" src="https://i.imgur.com/6OIPMR7.png">

#### Generating One time

<img style="max-width:600px" src="https://i.imgur.com/UC2Fi5W.png">
If onetime origin is invalid it will show error

<img style="max-width:600px" src="https://i.imgur.com/bMAYiCI.png">

## API

Server has 4 endpoints
Which are defined in [router.go](https://github.com/gynshu-one/goph-keeper/blob/main/server/api/router/router.go)
### /user/create
Creates new user with username and password from url params
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
Synchronizes user data with server. Server checks if user has session cookie via
[middleware.go](https://github.com/gynshu-one/goph-keeper/blob/main/server/api/middlewares/middleware.go)
which uses [session.go](https://github.com/gynshu-one/goph-keeper/blob/main/server/api/auth/session.go)
```
https://localhost:8080/user/sync
```
endpoint expects `POST` request with slice of `DataWrapper` structs in json body 
```go
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

struct in [general.go](https://github.com/gynshu-one/goph-keeper/blob/main/common/models/general.go)


`Sync` happens imminently after logging in and every item creation, deletion or update.

`Client` does not have cache of the data, it only stores username in ~/.goph-keeper/user.txt file.


Deleted items would force server to remove `DataWrapper`'s Data field and set `DeletedAt` field.
