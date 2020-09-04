package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
Before running this app, launch 'cqlsh' and execute:
CREATE KEYSPACE message_sender_rest_api with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE message_sender_rest_api.messages(id UUID, email text, title text, content text, magic_number int, created_at timestamp, PRIMARY KEY(id));
CREATE INDEX ON message_sender_rest_api.messages(magic_number);
CREATE INDEX ON message_sender_rest_api.messages(email);
*/

func main() {
	config := LoadConfig()

	// new db session
	session := NewCassandraSession()
	defer Close(session)
	// new message service
	ms := NewMessageService(session)
	// new message controller
	messageC := NewMessageController(ms, config.Mailgun)

	// router
	r := mux.NewRouter()
	r.HandleFunc("/api/message", messageC.WriteMessage).Methods("POST")
	r.HandleFunc("/api/send", messageC.SendMessages).Methods("POST")
	r.HandleFunc("/api/messages/{email}", messageC.GetMessagesByEmail).Methods("GET")

	http.ListenAndServe(":8080", r)
}
