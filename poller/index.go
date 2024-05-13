package poller

import (
	"errors"
	"time"

	"github.com/go-co-op/gocron/v2"
	"usepolymer.co/background/logger"
)


func BeginCronPoll(interval int, action func()) {
    s, err := gocron.NewScheduler()
	if err != nil {
        logger.Error(errors.New("error creating new gocron scheduler"), logger.LoggerOptions{
            Key: "error",
            Data: err,
        })
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
            time.Minute * time.Duration(interval),
            ),
            gocron.NewTask(action),
	)
	if err != nil {
        logger.Error(errors.New("error creating new gocron job"), logger.LoggerOptions{
            Key: "error",
            Data: err,
        })
	}
	// start the scheduler
	s.Start()
}