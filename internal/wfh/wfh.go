package wfh

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/roman-kulish/wfh/internal/slack"
)

const (
	command     = "/wfh"
	msgToday    = "<@%s> is working from home today"
	msgTomorrow = "<@%s> will be working from home tomorrow"
	msgMonday   = "<@%s> will be working from home on Monday"
	imgTitle    = "My excuse is ..."
)

type CommandHandler struct {
	Location       *time.Location
	ImageBaseURL   string
	NumberOfImages uint
}

func New(timezone, imageBaseURL string, numberOfImages uint) (*CommandHandler, error) {
	var location *time.Location
	var err error

	if timezone != "" {
		location, err = time.LoadLocation(timezone)
		if err != nil {
			return nil, err
		}
	} else {
		location = time.Local
	}

	if imageBaseURL != "" {
		if !strings.HasSuffix(imageBaseURL, "/") {
			imageBaseURL += "/"
		}

		imageBaseURL += "%x.jpg"
	}

	return &CommandHandler{Location: location, ImageBaseURL: imageBaseURL, NumberOfImages: numberOfImages}, nil
}

func (wfh *CommandHandler) Handle(req slack.CommandRequest) (slack.CommandResponse, error) {
	var msg string

	if req.Command != command {
		return slack.CommandResponse{}, fmt.Errorf("invalid command %s", req.Command)
	}

	now := time.Now().In(wfh.Location)
	then := time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, wfh.Location)

	switch {
	case now.Weekday() == 0 || now.Weekday() > 5 || (now.Weekday() == 5 && now.After(then)):
		msg = msgMonday
	case now.After(then):
		msg = msgTomorrow
	default:
		msg = msgToday
	}

	msg = fmt.Sprintf(msg, req.UserID)

	if txt := strings.TrimSpace(req.Text); txt != "" {
		txt = strings.Trim(txt, "_~*`")
		txt = "_" + txt + "_"

		msg = msg + ": " + txt // echo text back
	}

	res := slack.NewInChannelCommandResponse(msg)

	if wfh.ImageBaseURL != "" && wfh.NumberOfImages > 0 {
		rand.Seed(time.Now().UTC().UnixNano())

		index := rand.Intn(int(wfh.NumberOfImages))
		if index == 0 {
			index++
		}

		hash := sha1.Sum([]byte(strconv.Itoa(index)))

		res.AddAttachment(slack.Attachment{
			Title:    imgTitle,
			ImageURL: fmt.Sprintf(wfh.ImageBaseURL, hash),
		})
	}

	return res, nil
}
