package receiver

import (
	"errors"

	"github.com/gorilla/websocket"
	"github.com/segurosfalabella/imperium-worker/connection"
	log "github.com/sirupsen/logrus"
)

const passwordForSend = "alohomora"
const passwordForValidate = "imperio"

// JobProcessor interface
type JobProcessor interface {
	Execute() error
	FromJSON(text string)
	ToJSON() string
}

// Start function
func Start(conn connection.WsConn, jobProcessor JobProcessor) {
	err := auth(conn)
	if err != nil {
		log.Error(err.Error())
		return
	}
	loop(conn, jobProcessor)
}

func auth(conn connection.WsConn) error {
	conn.WriteMessage(websocket.TextMessage, []byte(passwordForSend))
	_, message, _ := conn.ReadMessage()
	if string(message) != passwordForValidate {
		return errors.New("server unknown")
	}
	log.Info("log in succeed!")
	return nil
}

func loop(conn connection.WsConn, jobProcessor JobProcessor) {
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("error reading messages: ", err)
			return
		}

		if err := parseJob(messageType, message, jobProcessor); err == nil {
			process(conn, jobProcessor)
		}
	}
}

func parseJob(messageType int, message []byte, jobProcessor JobProcessor) error {
	switch messageType {
	case websocket.TextMessage:
		jobProcessor.FromJSON(string(message))
	default:
		return errors.New("not supported format")
	}

	return nil
}

func process(conn connection.WsConn, jobProcessor JobProcessor) {
	if err := jobProcessor.Execute(); err != nil {
		log.Error(jobProcessor.ToJSON())
	}

	log.Info("process message succeed: ", jobProcessor.ToJSON())
	conn.WriteMessage(websocket.TextMessage, []byte(jobProcessor.ToJSON()))
}
