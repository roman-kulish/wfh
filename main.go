package main

import (
	"net/http"
	"time"

	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"encoding/hex"

	"github.com/roman-kulish/wfh/slash"
)

const (
	command     = "/wfh"
	msgToday    = "<!here> <@%s> is working from home today"
	msgTomorrow = "<!here> <@%s> will be working from home tomorrow"
	msgMonday   = "<!here> <@%s> will be working from home on Monday"
	imgTitle    = "My excuse is ..."
	port        = 8080
)

var (
	bucket string
	addr   string
)

type slashCommandHandler struct {
	Handler func(cmd slash.CommandRequest) (slash.CommandResponse, error)
	http.HandlerFunc
}

func (sc slashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if req.Command != command {
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	mux := *http.NewServeMux()

	mux.Handle(command, slashCommandHandler{
		Handler: func(req slash.CommandRequest) (slash.CommandResponse, error) {
			var message string

			loc, _ := time.LoadLocation("Australia/Sydney")
			now := time.Now().In(loc)
			then := time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, loc)

			switch true {
			case now.Weekday() == 0 || now.Weekday() > 5 || (now.Weekday() == 5 && now.After(then)):
				message = msgMonday
			case now.After(then):
				message = msgTomorrow
			default:
				message = msgToday
			}

			message = fmt.Sprintf(message, req.UserId)
			index := rand.Intn(250)

			if index == 0 {
				index++
			}

			hash := sha1.Sum([]byte(strconv.Itoa(index)))
			hashed := hex.EncodeToString(hash[:])

			res := slash.NewInChannelCommandResponse(message)

			res.AddAttachment(slash.Attachment{
				Title:    imgTitle,
				ImageUrl: fmt.Sprintf(bucket, hashed),
			})

			return res, nil
		},
	})

	server := http.Server{
		Addr:         addr,
		Handler:      &mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 3,
		IdleTimeout:  time.Second * 10,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func init() {
	if listen, ok := os.LookupEnv("APP_PORT"); ok {
		addr = fmt.Sprintf(":%s", listen)
	} else {
		addr = fmt.Sprintf(":%d", port)
	}
}
