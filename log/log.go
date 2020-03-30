// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package log has the interface for loggers in the sdk
package log

import (
	"fmt"
	"log"

	"github.com/cuttle-ai/octopus/lsp/config"
)

//Log interface to be implmented to be considered as a logger for web application
type Log interface {
	//Info logs the info
	Info(l ...interface{})
	//Debug logs debug specific information which will be skipped in production
	Debug(l ...interface{})
	//Warn logs warning info
	Warn(l ...interface{})
	//Error logs error info
	Error(l ...interface{})
	//Fatal logs the error and exits
	Fatal(l ...interface{})
}

//Log types for logger
const (
	//INFO is for informative logs
	INFO = "INFO"
	//DEBUG is for debugging the app
	DEBUG = "DEBUG"
	//WARN is for warning signatures
	WARN = "WARN"
	//ERROR is for errors
	ERROR = "ERROR"
	//PANIC is for panic log prefix
	PANIC = "PANIC"
)

//Logger which implements the Log
type Logger struct{}

//NewLogger returns the logger instance
func NewLogger() Logger {
	return Logger{}
}

//Info logs the info logs of the application
func (lg Logger) Info(l ...interface{}) {
	log.Print(INFO+": ", fmt.Sprintln(l...))
}

//Debug logs the debug logs of the application if debug logs are not switched off
func (lg Logger) Debug(l ...interface{}) {
	//Checking if Debug log is off
	if config.PRODUCTION == 0 {
		return
	}
	log.Print(DEBUG+": ", fmt.Sprintln(l...))
}

//Warn logs the warning logs of the application
func (lg Logger) Warn(l ...interface{}) {
	log.Print(WARN+": ", fmt.Sprintln(l...))
}

//Error logs the error logs of the application
func (lg Logger) Error(l ...interface{}) {
	log.Print(ERROR+": ", fmt.Sprintln(l...))
}

//Fatal is used to print logs for events which causes the app to exit
func (lg Logger) Fatal(l ...interface{}) {
	/*
	 * We will call log.Fatal
	 */
	log.Fatal(PANIC+": ", fmt.Sprintln(l...))
}
