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
	dT := interpreter.DataTypeString
	aggFn := interpreter.AggregationFnCount
	mes := false
	dim := false
	name := ""
	word := ""
	for _, v := range n.Metadata {
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
		}
	}
	return interpreter.ColumnNode{
		UID:           n.UID.String(),
		Word:          word,
		PUID:          n.PUID.String(),
		Name:          name,
		Children:      []interpreter.ValueNode{},
		Dimension:     dim,
		Measure:       mes,
		AggregationFn: aggFn,
		DataType:      dT,
	}
}

//FromColumn converts the interpreter column node to node
func (n Node) FromColumn(c interpreter.ColumnNode) Node {
	metadata := []NodeMetadata{}
	for _, v := range n.Metadata {
		metadata = append(metadata, v)
	}
	for i := 0; i < len(metadata); i++ {
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = c.Word
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
		}
	}
	return Node{}
}
