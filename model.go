package main

import (
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

var pageSize = 32

var errInvalidEmailFormat error = errors.New("invalid email format")
var errInvalidID error = errors.New("invalid ID")

// Message represents message data in the database
type Message struct {
	ID          gocql.UUID `json:"-"`
	Email       string     `json:"email"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	MagicNumber int        `json:"magic_number"`
	CreatedAt   time.Time  `json:"-"`
}

// MessageDB is an interface for interacting with message data in the database
type MessageDB interface {
	FindByEmail(email string) ([]Message, error)
	FindByMagicNumber(magic int) ([]Message, error)

	Create(message *Message) error
	Delete(id gocql.UUID) error
}

// MessageService is an interface for interacting with message model
type MessageService interface {
	MessageDB
	DeleteOldMessages(messages []Message) ([]Message, error)
}

type messageService struct {
	Session    *gocql.Session
	EmailRegex *regexp.Regexp
}

// NewMessageService creates MessageService instance
func NewMessageService(session *gocql.Session) MessageService {
	return &messageService{
		Session:    session,
		EmailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
	}
}

func (ms *messageService) FindByEmail(email string) ([]Message, error) {
	// validations
	email = strings.ToLower(email)
	if !ms.EmailRegex.MatchString(email) {
		return nil, errInvalidEmailFormat
	}

	// query to the database
	messages := make([]Message, 0)
	m := map[string]interface{}{}

	iter := ms.Session.Query("SELECT id, email, title, content, magic_number, created_at FROM messages WHERE email = ?", email).PageSize(pageSize).Iter()
	for iter.MapScan(m) {
		messages = append(messages, Message{
			ID:          m["id"].(gocql.UUID),
			Email:       m["email"].(string),
			Title:       m["title"].(string),
			Content:     m["content"].(string),
			MagicNumber: m["magic_number"].(int),
			CreatedAt:   m["created_at"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()
	if err != nil {
		log.Println("Failure during database query")
		return nil, err
	}

	return messages, nil
}

func (ms *messageService) FindByMagicNumber(magic int) ([]Message, error) {
	// validations

	// query to the database
	messages := make([]Message, 0)
	m := map[string]interface{}{}

	iter := ms.Session.Query("SELECT id, email, title, content, magic_number, created_at FROM messages WHERE magic_number = ?", magic).PageSize(pageSize).Iter()
	for iter.MapScan(m) {
		messages = append(messages, Message{
			ID:          m["id"].(gocql.UUID),
			Email:       m["email"].(string),
			Title:       m["title"].(string),
			Content:     m["content"].(string),
			MagicNumber: m["magic_number"].(int),
			CreatedAt:   m["created_at"].(time.Time),
		})
		m = map[string]interface{}{}
	}

	err := iter.Close()
	if err != nil {
		log.Println("Failure during database query")
		return nil, err
	}

	return messages, nil
}

func (ms *messageService) Create(message *Message) error {
	// validations
	message.Email = strings.ToLower(message.Email)
	if !ms.EmailRegex.MatchString(message.Email) {
		return errInvalidEmailFormat
	}

	// insertion into database
	err := ms.Session.Query("INSERT INTO messages(id, email, title, content, magic_number, created_at) VALUES(?,?,?,?,?,?)",
		gocql.TimeUUID(), message.Email, message.Title, message.Content, message.MagicNumber, time.Now()).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (ms *messageService) Delete(id gocql.UUID) error {
	// validations
	var nilID gocql.UUID
	if id == nilID {
		return errInvalidID
	}

	// delete from the database
	err := ms.Session.Query("DELETE FROM messages where id = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (ms *messageService) DeleteOldMessages(messages []Message) ([]Message, error) {
	retMessages := make([]Message, 0)
	fiveMinutes := time.Minute * 5

	for _, message := range messages {
		messageDuration := time.Now().Sub(message.CreatedAt)
		if messageDuration > fiveMinutes {
			err := ms.Delete(message.ID)
			if err != nil {
				log.Println("Failed to delete a record")
				return messages, err
			}
		} else {
			retMessages = append(retMessages, message)
		}
	}

	return retMessages, nil
}
