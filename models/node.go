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

const (
	//NodeMetadataPropWord is the metadata property of a node for word
	NodeMetadataPropWord = "Word"
	//NodeMetadataPropName is the metadata property of a node for name
	NodeMetadataPropName = "Name"
	//NodeMetadataPropDimension is the metadata property of a node for dimension
	NodeMetadataDimension = "Dimension"
	//NodeMetadataPropMeasure is the metadata property of a node for measure
	NodeMetadataMeasure = "Measure"
	//NodeMetadataPropAggregationFn is the metadata property of a node for aggregation function
	NodeMetadataAggregationFn = "AggregationFn"
	//NodeMetadataPropDataType is the metadata property of a node for data type
	NodeMetadataDataType = "DataType"
	//NodeMetadataPropDescription is the metadata property of a node for description
	NodeMetadataPropDescription = "Description"
)

const (
	//NodeMetadataPropValueTrue is the value to be put for true as metadata value
	NodeMetadataPropValueTrue = "true"
	//NodeMetadataPropValueFalse is the value to be put for false as metadata value
	NodeMetadataPropValueFalse = "false"
)

var (
	//NodeMetadataAggregationFns is the map containing the supported aggregation functions
	NodeMetadataAggregationFns = map[string]struct{}{
		interpreter.AggregationFnAvg:   struct{}{},
		interpreter.AggregationFnCount: struct{}{},
		interpreter.AggregationFnSum:   struct{}{},
	}
	//NodeMetadataDataTypes is the map containing the supported datatypes
	NodeMetadataDataTypes = map[string]struct{}{
		interpreter.DataTypeDate:   struct{}{},
		interpreter.DataTypeFloat:  struct{}{},
		interpreter.DataTypeInt:    struct{}{},
		interpreter.DataTypeString: struct{}{},
	}
)

//Node represents a octopus node's db record
type Node struct {
	gorm.Model
	//UID is the unique id of the node
	UID uuid.UUID
	//Type of the node
	Type interpreter.Type
	//PUID is the unique id of the parent node
	PUID uuid.UUID
	//DatasetID is the id of the dataset to which the node belongs to
	DatasetID uint
	//NodeMetadatas holds the metadata corresponding to the node
	NodeMetadatas []NodeMetadata
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
	dT := interpreter.DataTypeString
	aggFn := interpreter.AggregationFnCount
	mes := false
	dim := false
	name := ""
	word := ""
	description := ""
	for _, v := range n.NodeMetadatas {
		if v.Prop == NodeMetadataPropWord {
			word = v.Value
		} else if v.Prop == NodeMetadataPropName {
			name = v.Value
		} else if v.Prop == NodeMetadataDimension && v.Value == NodeMetadataPropValueTrue {
			dim = true
		} else if v.Prop == NodeMetadataMeasure && v.Value == NodeMetadataPropValueTrue {
			mes = true
		} else if v.Prop == NodeMetadataAggregationFn {
			if _, ok := NodeMetadataAggregationFns[v.Value]; ok {
				aggFn = v.Value
			}
		} else if v.Prop == NodeMetadataDataType {
			if _, ok := NodeMetadataDataTypes[v.Value]; ok {
				dT = v.Value
			}
		} else if v.Prop == NodeMetadataPropDescription {
			description = v.Value
		}
	}
	return interpreter.ColumnNode{
		UID:           n.UID.String(),
		Word:          []rune(word),
		PUID:          n.PUID.String(),
		Name:          name,
		Children:      []interpreter.ValueNode{},
		Dimension:     dim,
		Measure:       mes,
		AggregationFn: aggFn,
		DataType:      dT,
		Description:   description,
	}
}

//FromColumn converts the interpreter column node to node
func (n Node) FromColumn(c interpreter.ColumnNode) Node {
	metadata := []NodeMetadata{}
	for _, v := range n.NodeMetadatas {
		metadata = append(metadata, v)
	}
	if len(metadata) == 0 {
		metadata = append(metadata, NodeMetadata{
			Prop: NodeMetadataPropWord,
		}, NodeMetadata{
			Prop: NodeMetadataPropName,
		}, NodeMetadata{
			Prop: NodeMetadataDimension,
		}, NodeMetadata{
			Prop: NodeMetadataMeasure,
		}, NodeMetadata{
			Prop: NodeMetadataAggregationFn,
		}, NodeMetadata{
			Prop: NodeMetadataDataType,
		}, NodeMetadata{
			Prop: NodeMetadataPropDescription,
		})
	}
	for i := 0; i < len(metadata); i++ {
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = string(c.Word)
		} else if metadata[i].Prop == NodeMetadataPropName {
			metadata[i].Value = c.Name
		} else if metadata[i].Prop == NodeMetadataDimension {
			if c.Dimension {
				metadata[i].Value = NodeMetadataPropValueTrue
			} else {
				metadata[i].Value = NodeMetadataPropValueFalse
			}
		} else if metadata[i].Prop == NodeMetadataMeasure {
			if c.Measure {
				metadata[i].Value = NodeMetadataPropValueTrue
			} else {
				metadata[i].Value = NodeMetadataPropValueFalse
			}
		} else if metadata[i].Prop == NodeMetadataAggregationFn {
			if _, ok := NodeMetadataAggregationFns[c.AggregationFn]; ok {
				metadata[i].Value = c.AggregationFn
			} else {
				metadata[i].Value = interpreter.AggregationFnCount
			}
		} else if metadata[i].Prop == NodeMetadataDataType {
			if _, ok := NodeMetadataDataTypes[c.DataType]; ok {
				metadata[i].Value = c.DataType
			} else {
				metadata[i].Value = interpreter.DataTypeString
			}
		} else if metadata[i].Prop == NodeMetadataPropDescription {
			metadata[i].Value = c.Description
		}
	}
	uid, _ := uuid.Parse(c.UID)
	puid, _ := uuid.Parse(c.PUID)
	return Node{
		Model:         n.Model,
		UID:           uid,
		Type:          c.Type(),
		PUID:          puid,
		DatasetID:     n.DatasetID,
		NodeMetadatas: metadata,
	}
}
