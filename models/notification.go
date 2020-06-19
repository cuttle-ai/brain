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

const (
	//InfoNotification is the notitfication for info type of notification
	InfoNotification = "INFO_NOTIFICATION"
)
