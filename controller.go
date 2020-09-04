package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go"
)

// MessageController is a controller struct responsible for handling message resources
type MessageController struct {
	ms MessageService
	mg mailgun.Mailgun
}

// NewMessageController creates new message controller
func NewMessageController(ms MessageService, mg MailgunConfig) *MessageController {
	return &MessageController{
		ms: ms,
		mg: mailgun.NewMailgun(mg.Domain, mg.APIKey),
	}
}

// WriteMessage handles POST /api/message
func (mc *MessageController) WriteMessage(w http.ResponseWriter, r *http.Request) {
	var message Message

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read user's input")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println("Failed to parse user's input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = mc.ms.Create(&message)
	if err != nil {
		log.Println("Failed to create new record")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// SendMessages handles POST /api/send
func (mc *MessageController) SendMessages(w http.ResponseWriter, r *http.Request) {
	var message Message

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Failed to read user's input")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &message)
	if err != nil {
		log.Println("Failed to parse user's input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	messages, err := mc.ms.FindByMagicNumber(message.MagicNumber)
	if err != nil {
		log.Println("Failed to query for the messages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messages, err = mc.ms.DeleteOldMessages(messages)
	if err != nil {
		log.Println("Failed to delete old messages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, message := range messages {
		// send email
		err = mc.sendEmail(message)
		if err != nil {
			log.Println("Failed to send the message")
			log.Println(err)
		}
		// delete from the database
		err = mc.ms.Delete(message.ID)
		if err != nil {
			log.Println("Failed to delete the message")
		}
	}
}

// GetMessagesByEmail handles GET /api/messages/{email}
func (mc *MessageController) GetMessagesByEmail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	email, prs := params["email"]
	if prs == false {
		log.Println("Failed to parse url")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messages, err := mc.ms.FindByEmail(email)
	if err != nil {
		log.Println("Failed to query for the messages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	messages, err = mc.ms.DeleteOldMessages(messages)
	if err != nil {
		log.Println("Failed to delete old messages")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.MarshalIndent(messages, "", " ")
	if err != nil {
		log.Println("Failed to prepare response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, string(response))
}

func (mc *MessageController) sendEmail(message Message) error {
	m := mc.mg.NewMessage(
		"pduzinki's message sender rest api <mailgun@pduzinki.com>",
		message.Title,
		message.Content,
		message.Email,
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, _, err := mc.mg.Send(ctx, m)
	return err
}
