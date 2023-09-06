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
sudo chmod +x server/cert/gen_keys.sh
cd server/cert && ./gen_keys.sh
# runs mongodb
docker-compose up --build -d
# builds the server
go build -o ./server_cmd ./server/cmd/main.go
# runs the server
./server_cmd -c ./server/config.json
```
This will generate TSL certs build docker image and run it on port 8080 as well as mongodb on port 27017

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


see [Makefile](https://github.com/gynshu-one/goph-keeper/blob/main/Makefile)


## Signing up and logging in
<img style="max-width:600px" src="https://i.imgur.com/DDsrsM6.png">

For simplicity, the password and secret are stored in the OS keychain after registration.

Every other "log in" will grab username from  `/temp/session_id` file and check OS default `keychain` for password and secret.

Secret is used for encrypting local data 

After logging in to a server,client will receive session cookie for 24 hours.

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

Server compares received items to users previous 
items if presented, and updates according to `updated_at` Unix time field. If client sends empty body - response would be all stored user items on server side.


`Client` does not have cache of the data, it only stores session cookie and username in /temp/session_id file. 


Deleted items would force server to remove `DataWrapper`'s Data field and set `DeletedAt` field.
