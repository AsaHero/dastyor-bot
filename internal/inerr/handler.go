package inerr

import (
	"fmt"

	"github.com/AsaHero/dastyor-bot/pkg/logger"
	"github.com/AsaHero/dastyor-bot/pkg/utility"
	"github.com/sirupsen/logrus"
)

func Err(err error) error {
	if err == nil {
		err = fmt.Errorf("unknown error")
	}

	scope, caller, callee, location := utility.GetFrameData(2)

	logger.Error(err.Error(), logrus.Fields{
		"scope":    scope,
		"caller":   caller,
		"callee":   callee,
		"location": location,
	})

	return err
}

func Newf(format string, msg ...any) error {
	scope, caller, callee, location := utility.GetFrameData(2)

	err := fmt.Errorf(format, msg...)

	logger.Error(err.Error(), logrus.Fields{
		"scope":    scope,
		"caller":   caller,
		"callee":   callee,
		"location": location,
	})

	return err
}
