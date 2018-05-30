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
	msgToday    = "<!here> <@%s> is working from home today"
	msgTomorrow = "<!here> <@%s> will be working from home tomorrow"
	msgMonday   = "<!here> <@%s> will be working from home on Monday"
	imgTitle    = "My excuse is %s"
	// epsilon     = "..."
)

type wfh struct {
	Location       *time.Location
	ImageBaseUrl   string
	NumberOfImages uint
}

func New(timezone, imageBaseUrl string, numberOfImages uint) (*wfh, error) {
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

	if imageBaseUrl != "" {
		if !strings.HasSuffix(imageBaseUrl, "/") {
			imageBaseUrl = imageBaseUrl + "/"
		}

		imageBaseUrl = imageBaseUrl + "%x.jpg"
	}

	return &wfh{Location: location, ImageBaseUrl: imageBaseUrl, NumberOfImages: numberOfImages}, nil
}

func (wfh *wfh) Handle(req slack.CommandRequest) (slack.CommandResponse, error) {
	var msg string

	if req.Command != command {
		return slack.CommandResponse{}, fmt.Errorf("invalid command %s", req.Command)
	}

	now := time.Now().In(wfh.Location)
	then := time.Date(now.Year(), now.Month(), now.Day(), 10, 15, 0, 0, wfh.Location)

	switch true {
	case now.Weekday() == 0 || now.Weekday() > 5 || (now.Weekday() == 5 && now.After(then)):
		msg = msgMonday
	case now.After(then):
		msg = msgTomorrow
	default:
		msg = msgToday
	}

	msg = fmt.Sprintf(msg, req.UserId)
	res := slack.NewInChannelCommandResponse(msg)

	if wfh.ImageBaseUrl != "" && wfh.NumberOfImages > 0 {
		rand.Seed(time.Now().UTC().UnixNano())

		index := rand.Intn(int(wfh.NumberOfImages))

		if index == 0 {
			index++
		}

		hash := sha1.Sum([]byte(strconv.Itoa(index)))

		res.AddAttachment(slack.Attachment{
			Title:    imgTitle,
			ImageUrl: fmt.Sprintf(wfh.ImageBaseUrl, hash),
		})
	}

	return res, nil
}
