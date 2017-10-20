package main

import (
	"net/http"
	"time"

	"bytes"
	"encoding/json"

	"crypto/sha1"
	"fmt"
	"math/rand"

	"strconv"

	"github.com/roman-kulish/wfh/slash"
)

const (
	endpoint    = "/wfh"
	bucket      = "https://storage.googleapis.com/wfh/%x.jpg"
	msgToday    = "@here <@%s> is working from today"
	msgTomorrow = "@here <@%s> will work from home tomorrow"
	msgMonday   = "@here <@%s> will work from home on Monday"
)

type SlashCommandHandler struct {
	Handler func(cmd slash.CommandRequest) (slash.CommandResponse, error)
	http.HandlerFunc
}

func (sc SlashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jsonData bytes.Buffer

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	req, err := slash.NewCommandRequest(r)

	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	res, err := sc.Handler(req)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if req.Command != endpoint {
		http.NotFound(w, r)
		return
	}

	if err := json.NewEncoder(&jsonData).Encode(&res); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData.Bytes())
}

func main() {
	handler := SlashCommandHandler{
		Handler: func(req slash.CommandRequest) (slash.CommandResponse, error) {
			var message = msgToday

			loc, _ := time.LoadLocation("Australia/Sydney")
			now := time.Now().In(loc)
			then := time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, loc)

			if now.After(then) {
				if now.Weekday() > 4 {
					message = msgMonday
				} else {
					message = msgTomorrow
				}
			}

			message = fmt.Sprintf(message, req.UserId)
			index := rand.Intn(250)

			if index == 0 {
				index++
			}

			hash := sha1.Sum([]byte(strconv.Itoa(index)))

			res := slash.NewInChannelCommandResponse(message)

			res.Attachments = append(res.Attachments, slash.Attachment{
				ImageUrl: fmt.Sprintf(bucket, hash),
			})

			return res, nil
		},
	}

	mux := *http.NewServeMux()
	mux.Handle(endpoint, handler)

	server := http.Server{
		Addr:         ":8080",
		Handler:      &mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 10,
	}

	server.ListenAndServe()
}
