FROM golang

WORKDIR /goph-keeper
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
RUN go mod verify
COPY ./server ./server
COPY ./common ./common
RUN go build -o /goph-keeper server/cmd/main.go  \
    && chmod +x main  \