// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"github.com/cuttle-ai/octopus/interpreter"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

/*
 * This file contains the model implementation of node's db model
 */

//Node represents a octopus node's db record
type Node struct {
	gorm.Model
	//UID is the unique id of the node
	UID uuid.UUID
	//Type of the node
	Type interpreter.Type
	//PID is the id of the parent of the node
	PID uint
	//DatasetID is the id of the dataset to which the node belongs to
	DatasetID uint
	//Metadata holds the metadata corresponding to the node
	Metadata []NodeMetadata
}

//NodeMetadata stores the metadata associated with a node
type NodeMetadata struct {
	gorm.Model
	//NodeID is the id of the node to which the metadata belongs to
	NodeID uint
	//Prop stores the metadata property
	Prop string
	//Value stores the metadata value
	Value string
}

//ColumnNode returns column node converted form of the node
func (n Node) ColumnNode() interpreter.ColumnNode {
	return interpreter.ColumnNode{}
}

//ColumnToNode converts the interpreter column node to node
func ColumnToNode(c interpreter.ColumnNode) Node {
	return Node{}
}
