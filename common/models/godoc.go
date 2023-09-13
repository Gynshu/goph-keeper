// Package models contains all the models used in the application
// Here are defined important methods for each model and BasicData (interface), DataWrapper (struct)
// BasicData is an interface for all data types and provides methods to encrypt and decrypt data
// DataWrapper is a struct that wraps BasicData and provides additional information about the data
// such as owner id, type, name, updated_at, created_at, deleted_at. This is useful to store data remotely knowing
// nothing about the data itself
// Login models also contains methods register and generate one-time password
package models
