package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roman-kulish/wfh/internal/slack"
	"github.com/roman-kulish/wfh/internal/wfh"
)

const timeout = 2 * time.Second

var (
	timezone       string
	imageBaseURL   string
	numberOfImages string
	command        *wfh.CommandHandler
	client         *http.Client

	ErrEmptyRequest = errors.New("HTTP POST body is empty")
)

func ParseBody(req string) (slack.CommandRequest, error) {
	values, err := url.ParseQuery(req)

	if err != nil {
		return slack.CommandRequest{}, err
	}

	return slack.CommandRequest{
		Token:          values.Get("token"),
		TeamID:         values.Get("team_id"),
		TeamDomain:     values.Get("team_domain"),
		EnterpriseID:   values.Get("enterprise_id"),
		EnterpriseName: values.Get("enterprise_name"),
		ChannelID:      values.Get("channel_id"),
		ChannelName:    values.Get("channel_name"),
		UserID:         values.Get("user_id"),
		UserName:       values.Get("user_name"),
		Command:        values.Get("command"),
		Text:           values.Get("text"),
		ResponseURL:    values.Get("response_url"),
		TriggerID:      values.Get("trigger_id"),
	}, nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrEmptyRequest
	}

	req, err := ParseBody(request.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	res, err := command.Handle(req)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	var jsonData bytes.Buffer
	encoder := json.NewEncoder(&jsonData)
	encoder.SetEscapeHTML(false)

	if err = encoder.Encode(&res); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	mRes, err := client.Post(req.ResponseURL, "application/json", &jsonData)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	io.Copy(io.Discard, mRes.Body)
	mRes.Body.Close()

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

func main() {
	timezone = os.Getenv("WFH_TIMEZONE")
	imageBaseURL = os.Getenv("WFH_IMAGE_BASE_URL")
	numberOfImages = os.Getenv("WFH_NUMBER_OF_IMAGES")

	num, err := strconv.Atoi(numberOfImages)
	if err != nil {
		panic(err)
	} else if num <= 0 {
		panic("number of images must be a positive integer")
	}

	command, err = wfh.New(timezone, imageBaseURL, uint(num))
	if err != nil {
		panic(err)
	}

	client = &http.Client{Timeout: timeout}

	lambda.Start(Handler)
}
