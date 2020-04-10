// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package models contains the models required by the brain platform
package models

import (
	"github.com/jinzhu/gorm"
)

const (
	//DatasetSourceFile indicates that the dataset source is file
	DatasetSourceFile = "FILE"
)

//Dataset represents a dataset
type Dataset struct {
	gorm.Model
	//Name of the dataset
	Name string
	//Description is the description for the dataset
	Description string
	//UserID is the id of the user with whom the file is associated with
	UserID uint
	//Source is the type of dataset source. It can be file, database etc
	Source string
	//ResourceID is the lid of the underlying dataset like file id for a dataset who source is file
	ResourceID uint
	//UploadedDataset is the uploaded data set info
	UploadedDataset interface{} `gorm:"-"`
	//TableCreated indicates whether the table is created for the dataset in the datastore
	TableCreated bool
	//DatastoreID is the id of the datastore where the data is physically stored for the dataset
	DatastoreID uint
}

const (
	//DatasetAccessTypeDashboard gives minimum access to the user.
	//The user won't get the dataset listed in datasets list. But will have minimum access to see the data through dashboard
	DatasetAccessTypeDashboard = 0
	//DatasetAccessTypeCreator gives user access to delete/update and all the previleges on the dataset
	DatasetAccessTypeCreator = 10
)

//DatsetUserMapping has the mapping of a dataset to  user.
//this includes the access type. All users with creator access and dashboard access will be listed in this table.
type DatsetUserMapping struct {
	gorm.Model
	//DatasetID is the ID of the dataset
	DatasetID uint
	//UserID is the ID of the user
	UserID uint
	//AccessType is the type of access for the user to the dataset
	AccessType int
}
