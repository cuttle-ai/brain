// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

//Notification is the data translation object for sending notifications
type Notification struct {
	//Event is the event to be called
	Event string
	//Payload is the payload to be send with the notification
	Payload interface{}
}

//ActionNotificationPayload has info for sending an action notification with message and required action
type ActionNotificationPayload struct {
	//Message of the notification
	Message string `json:"message,omitempty"`
	//Action to be performed
	Action string `json:"action,omitempty"`
}

const (
	//ActionFetchDatasets directs the client to refetch the list datasets api
	ActionFetchDatasets = "DATASETS"
)

const (
	//InfoNotification is the notitfication for info type of notification
	InfoNotification = "INFO_NOTIFICATION"
	//ErrorNotification is the notitfication for error type of notification
	ErrorNotification = "ERROR_NOTIFICATION"
	//SuccessNotification is the notitfication for success type of notification
	SuccessNotification = "SUCCESS_NOTIFICATION"
	//ActionNotification is the notification that will show a message and suggests to do an action.
	//eg:- When a dataset is deleted, user has to be notified on the same and the existing datasets list must be updated
	ActionNotification = "ACTION_NOTIFICATION"
)
