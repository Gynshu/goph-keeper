// Package storage provides a simple in-memory storage for the goph-keeper service.
// Every  operation is performed on the cache and asynchronously on the database.
// Storage itself, does not know Which exact model it is storing, it only knows that it is storing models.UserData.
// Sensitive fields are encrypted by client before sending to the server.
package storage
