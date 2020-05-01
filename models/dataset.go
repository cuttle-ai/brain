// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package models contains the models required by the brain platform
package models

import (
	"github.com/cuttle-ai/brain/log"
	"github.com/cuttle-ai/octopus/interpreter"
	"github.com/google/uuid"
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

//GetColumns get the columns corresponding to a dataset
func (d Dataset) GetColumns(conn *gorm.DB) ([]Node, error) {
	result := []Node{}
	err := conn.Set("gorm:auto_preload", true).Where("dataset_id = ? and type = ?", d.ID, interpreter.Column).Find(&result).Error
	return result, err
}

//GetTable get the tables corresponding to a dataset
func (d Dataset) GetTable(conn *gorm.DB) (Node, error) {
	result := []Node{}
	err := conn.Set("gorm:auto_preload", true).Where("dataset_id = ? and type = ?", d.ID, interpreter.Table).Find(&result).Error
	if len(result) > 0 {
		return result[0], nil
	}
	return Node{}, err
}

//Get will find the dataset values and set in the instance. Returns an error if couldn't find
func (d *Dataset) Get(conn *gorm.DB) error {
	return conn.Where("user_id = ? and id = ?", d.UserID, d.ID).Find(d).Error
}

//UpdateColumns updates the columns in the database. It will create the columns if not existing
func (d *Dataset) UpdateColumns(l log.Log, conn *gorm.DB, cols []Node) ([]Node, error) {
	/*
	 * We will use the db transactions to start update
	 * If id exists we will update
	 * else we will create the model
	 */
	//starting the transaction
	tx := conn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return nil, err
	}

	//will iterate through the cols for create/update
	for i := 0; i < len(cols); i++ {
		cols[i].DatasetID = d.ID
		//if id doesn't exists we will create the node
		if cols[i].ID == 0 {
			cols[i].UID = uuid.New()
			err := tx.Create(&cols[i]).Error
			if err != nil {
				l.Error("error while creating the column node for", cols[i].DatasetID, "at index", i)
				tx.Rollback()
				return nil, err
			}
			continue
		}
		//else we will update the node
		for j := 0; j < len(cols[i].NodeMetadatas); j++ {
			err := tx.Save(&(cols[i].NodeMetadatas[j])).Error
			if err != nil {
				l.Error("error while updating metadata of the column node for", cols[i].ID, cols[i].NodeMetadatas[j].Prop, cols[i].NodeMetadatas[j].ID)
				tx.Rollback()
				return nil, err
			}
		}
	}
	return cols, tx.Commit().Error
}
