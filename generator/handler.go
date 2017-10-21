package generator

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"math/rand"

	"github.com/roman-kulish/wfh/slash"
)

const (
	Command     = "/wfh"
	bucket      = "https://storage.googleapis.com/wfh/%x.jpg"
	msgToday    = "<!here> <@%s> is working from home today"
	msgTomorrow = "<!here> <@%s> will work from home tomorrow"
	msgMonday   = "<!here> <@%s> will work from home on Monday"
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

	if req.Command != Command {
		http.NotFound(w, r)
		return
	}

	encoder := json.NewEncoder(&jsonData)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(&res); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData.Bytes())
}

func NewHandler() SlashCommandHandler {
	rand.Seed(time.Now().UTC().UnixNano())

	return SlashCommandHandler{
		Handler: func(req slash.CommandRequest) (slash.CommandResponse, error) {
			var message = msgToday

			loc, _ := time.LoadLocation("Australia/Sydney")
			now := time.Now().In(loc)
			then := time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, loc)

			if now.Weekday() > 4 {
				message = msgMonday
			} else if now.After(then) {
				message = msgTomorrow
			}

			message = fmt.Sprintf(message, req.UserId)
			index := rand.Intn(250)

			if index == 0 {
				index++
			}

			hash := sha1.Sum([]byte(strconv.Itoa(index)))

			res := slash.NewInChannelCommandResponse(message)

			res.AddAttachment(slash.Attachment{
				Title:    "Excuse",
				ImageUrl: fmt.Sprintf(bucket, hash),
			})

			return res, nil
		},
	}
}
