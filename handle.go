package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
)

type handler struct {
	chanAccessToken string
}

func newHandler(chanAccessToken string) *handler {
	return &handler{chanAccessToken}
}

func (h *handler) Handle(ctx echo.Context) error {
	req := new(lineRequest)
	if err := ctx.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Print(req)
	h.Reply(req.Events[0].ReplyToken, req.Events[0].Message.Text)
	return ctx.JSON(http.StatusOK, nil)
}

const replyURL = "https://api.line.me/v2/bot/message/reply"

func (h *handler) Reply(replyToken, text string) error {
	body, err := json.Marshal(newLineResponse(replyToken, text))
	if err != nil {
		return fmt.Errorf("failed to marshal request; %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, replyURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to make new request; %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+h.chanAccessToken)
	req.Header.Add("Content-Type", "application/json")

	// log request
	reqBytes, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Printf("failed to dump request; %v", err)
	} else {
		log.Print(string(reqBytes))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to post request; %w", err)
	}
	defer resp.Body.Close()

	// log response
	respBytes, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Printf("failed to dump response; %v", err)
	} else {
		log.Print(string(respBytes))
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("got non 2xx status code; %d; %s", resp.StatusCode, resp.Status)
	}

	return nil
}

type lineRequest struct {
	Events []event `json:"events"`
}

type event struct {
	Message    message `json:"message"`
	ReplyToken string  `json:"replyToken"`
}

type lineResponse struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []message `json:"messages"`
}

func newLineResponse(replyToken string, text string) lineResponse {
	return lineResponse{ReplyToken: replyToken, Messages: []message{newMessage(text)}}
}

type message struct {
	Typ  string `json:"type"`
	Text string `json:"text"`
}

func newMessage(text string) message {
	return message{Typ: "text", Text: text}
}
