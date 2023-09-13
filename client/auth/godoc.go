// Package auth is Custom auth system using session IDs
// I is about storing and retrieving passwords and secrets from OS keyring
// for simplicity. It also stores current user name in a file in user's home directory.
// Every time user starts the client, it reads the file and sets CurrentUser.Username.
// If the file does not exist, In this case user must login again.
package auth
