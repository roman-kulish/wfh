package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/roman-kulish/wfh/internal/slack"
	"github.com/roman-kulish/wfh/internal/wfh"
)

var (
	timezone       string
	imageBaseUrl   string
	numberOfImages string
	client         *http.Client
	wg             sync.WaitGroup
)

type slashCommandHandler struct {
	Handler func(cmd slack.CommandRequest) (slack.CommandResponse, error)
}

func (sch slashCommandHandler) GetCommandRequest(r *http.Request) (slack.CommandRequest, error) {
	if err := r.ParseForm(); err != nil {
		return slack.CommandRequest{}, err
	}

	return slack.CommandRequest{
		Token:          r.PostForm.Get("token"),
		TeamId:         r.PostForm.Get("team_id"),
		TeamDomain:     r.PostForm.Get("team_domain"),
		EnterpriseId:   r.PostForm.Get("enterprise_id"),
		EnterpriseName: r.PostForm.Get("enterprise_name"),
		ChannelId:      r.PostForm.Get("channel_id"),
		ChannelName:    r.PostForm.Get("channel_name"),
		UserId:         r.PostForm.Get("user_id"),
		UserName:       r.PostForm.Get("user_name"),
		Command:        r.PostForm.Get("command"),
		Text:           r.PostForm.Get("text"),
		ResponseUrl:    r.PostForm.Get("response_url"),
		TriggerId:      r.PostForm.Get("trigger_id"),
	}, nil
}

func (sch slashCommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var jsonData bytes.Buffer

	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	req, err := sch.GetCommandRequest(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res, err := sch.Handler(req)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(&jsonData)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(&res); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	wg.Add(1)

	go func() {
		client.Post(req.ResponseUrl, "application/json", &jsonData)

		wg.Done()
	}()
}

func main() {
	var addr string

	num, err := strconv.Atoi(numberOfImages)

	if err != nil {
		panic(err)
	} else if num <= 0 {
		panic("number of images must be a positive integer")
	}

	command, err := wfh.New(timezone, imageBaseUrl, uint(num))

	if err != nil {
		panic(err)
	}

	client = &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	mux := *http.NewServeMux()

	mux.Handle("/wfh", slashCommandHandler{
		Handler: func(cmd slack.CommandRequest) (slack.CommandResponse, error) {
			return command.Handle(cmd)
		},
	})

	if listen, ok := os.LookupEnv("WFH_PORT"); ok {
		addr = fmt.Sprintf(":%s", listen)
	} else {
		addr = ":8080"
	}

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

	wg.Wait()
}
