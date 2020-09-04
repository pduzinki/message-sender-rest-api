package main

import (
	"github.com/gocql/gocql"
)

// NewCassandraSession starts Cassandra DB session, panics on failure
func NewCassandraSession() *gocql.Session {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "message_sender_rest_api"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	return session
}

// Close closes given Cassanda DB session
func Close(session *gocql.Session) {
	session.Close()
}
