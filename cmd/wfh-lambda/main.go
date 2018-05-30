package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/roman-kulish/wfh/internal/slack"
	"github.com/roman-kulish/wfh/internal/wfh"
)

var (
	timezone       string
	imageBaseUrl   string
	numberOfImages string

	ErrEmptyRequest           = errors.New("HTTP POST body is empty")
	ErrNegativeNumberOfImages = errors.New("number of images must be a positive integer")
)

func ParseBody(req string) (slack.CommandRequest, error) {
	values, err := url.ParseQuery(req)

	if err != nil {
		return slack.CommandRequest{}, err
	}

	return slack.CommandRequest{
		Token:          values.Get("token"),
		TeamId:         values.Get("team_id"),
		TeamDomain:     values.Get("team_domain"),
		EnterpriseId:   values.Get("enterprise_id"),
		EnterpriseName: values.Get("enterprise_name"),
		ChannelId:      values.Get("channel_id"),
		ChannelName:    values.Get("channel_name"),
		UserId:         values.Get("user_id"),
		UserName:       values.Get("user_name"),
		Command:        values.Get("command"),
		Text:           values.Get("text"),
		ResponseUrl:    values.Get("response_url"),
		TriggerId:      values.Get("trigger_id"),
	}, nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var jsonData bytes.Buffer

	num, err := strconv.Atoi(numberOfImages)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	} else if num <= 0 {
		return events.APIGatewayProxyResponse{}, ErrNegativeNumberOfImages
	}

	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrEmptyRequest
	}

	req, err := ParseBody(request.Body)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	command, err := wfh.New(timezone, imageBaseUrl, uint(num))

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	res, err := command.Handle(req)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	encoder := json.NewEncoder(&jsonData)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(&res); err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       jsonData.String(),
		StatusCode: 200,
	}, nil
}

func main() {
	timezone = os.Getenv("WFH_TIMEZONE")
	imageBaseUrl = os.Getenv("WFH_IMAGE_BASE_URL")
	numberOfImages = os.Getenv("WFH_NUMBER_OF_IMAGES")

	lambda.Start(Handler)
}
