// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"strconv"

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
	NodeMetadataPropDimension = "Dimension"
	//NodeMetadataPropMeasure is the metadata property of a node for measure
	NodeMetadataPropMeasure = "Measure"
	//NodeMetadataPropAggregationFn is the metadata property of a node for aggregation function
	NodeMetadataPropAggregationFn = "AggregationFn"
	//NodeMetadataPropDataType is the metadata property of a node for data type
	NodeMetadataPropDataType = "DataType"
	//NodeMetadataPropDescription is the metadata property of a node for description
	NodeMetadataPropDescription = "Description"
	//NodeMetadataPropDefaultDateFieldUID is the metadata property of a node for default date field uid
	NodeMetadataPropDefaultDateFieldUID = "DefaultDateFieldUID"
	//NodeMetadataPropDatastoreID is the metadata property of a node for giving the datastore to which the node belongs to
	NodeMetadataPropDatastoreID = "NodeMetadataPropDatastoreID"
	//NodeMetadataPropKBType is the metadata property of a knowledge base node for giving kb type of the kb node
	NodeMetadataPropKBType = "KBType"
	//NodeMetadataPropOperation is the metadata property of a operation node for giving the type of the operation
	NodeMetadataPropOperation = "Operation"
	//NodeMetadataPropDateFormat is the metadata property of a column's csv data if the given column is of data type date
	NodeMetadataPropDateFormat = "DateFormat"
)

const (
	//NodeMetadataPropValueTrue is the value to be put for true as metadata value
	NodeMetadataPropValueTrue = "true"
	//NodeMetadataPropValueFalse is the value to be put for false as metadata value
	NodeMetadataPropValueFalse = "false"
	//NodeMetadataPropValueSystemKB is the value to be put for SystemKB as KBType metadata value
	NodeMetadataPropValueSystemKB = "1"
	//NodeMetadataPropValueUserKB is the value to be put for UserKB as KBType metadata value
	NodeMetadataPropValueUserKB = "2"
	//NodeMetadataPropValueEqOperator is the value to be put for Equal operator as operation of operator node
	NodeMetadataPropValueEqOperator = "="
	//NodeMetadataPropValueNotEqOperator is the value to be put for Not Equal operator as operation of operator node
	NodeMetadataPropValueNotEqOperator = "<>"
	//NodeMetadataPropValueGreaterOperator is the value to be put for Greater than or Equal operator as operation of operator node
	NodeMetadataPropValueGreaterOperator = ">="
	//NodeMetadataPropValueLessOperator is the value to be put for Less than or Equal operator as operation of operator node
	NodeMetadataPropValueLessOperator = "<="
	//NodeMetadataPropValueContainsOperator is the value to be put for Contains operator as operation of operator node
	NodeMetadataPropValueContainsOperator = "HAS"
	//NodeMetadataPropValueLikeOperator is the value to be put for Like operator as operation of operator node
	NodeMetadataPropValueLikeOperator = "LIKE"
)

var (
	//NodeMetadataAggregationFns is the map containing the supported aggregation functions
	NodeMetadataAggregationFns = map[string]struct{}{
		interpreter.AggregationFnAvg:   {},
		interpreter.AggregationFnCount: {},
		interpreter.AggregationFnSum:   {},
	}
	//NodeMetadataDataTypes is the map containing the supported datatypes
	NodeMetadataDataTypes = map[string]struct{}{
		interpreter.DataTypeDate:   {},
		interpreter.DataTypeFloat:  {},
		interpreter.DataTypeInt:    {},
		interpreter.DataTypeString: {},
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
	//Parent denotes the the parent for the node
	Parent *Node `gorm:"-"`
}

//NodeMetadata stores the metadata associated with a node
type NodeMetadata struct {
	gorm.Model
	//NodeID is the id of the node to which the metadata belongs to
	NodeID uint
	//DatasetID is the id of the dataset to which the node belongs to
	DatasetID uint
	//Prop stores the metadata property
	Prop string
	//Value stores the metadata value
	Value string
}

//InterpreterNode will convert a node to corresponding interpreter node
func (n Node) InterpreterNode() (interpreter.Node, bool) {
	switch n.Type {
	case interpreter.Column:
		cN := n.ColumnNode()
		return &cN, true
	case interpreter.Table:
		tN := n.TableNode()
		return &tN, true
	case interpreter.KnowledgeBase:
		kN := n.KnowledgeBaseNode()
		return &kN, true
	case interpreter.Operator:
		oN := n.OperatorNode()
		return &oN, true
	default:
		return nil, false
	}
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
	dateFormat := ""
	for _, v := range n.NodeMetadatas {
		if v.Prop == NodeMetadataPropWord {
			word = v.Value
		} else if v.Prop == NodeMetadataPropName {
			name = v.Value
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueTrue {
			dim = true
		} else if v.Prop == NodeMetadataPropMeasure && v.Value == NodeMetadataPropValueTrue {
			mes = true
		} else if v.Prop == NodeMetadataPropAggregationFn {
			if _, ok := NodeMetadataAggregationFns[v.Value]; ok {
				aggFn = v.Value
			}
		} else if v.Prop == NodeMetadataPropDataType {
			if _, ok := NodeMetadataDataTypes[v.Value]; ok {
				dT = v.Value
			}
		} else if v.Prop == NodeMetadataPropDescription {
			description = v.Value
		} else if v.Prop == NodeMetadataPropDateFormat {
			dateFormat = v.Value
		}
	}
	result := interpreter.ColumnNode{
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
		DateFormat:    dateFormat,
	}
	if n.Parent != nil && n.PUID.String() == n.Parent.UID.String() && n.Parent.Type == interpreter.Table {
		pN := n.Parent.TableNode()
		result.PN = &pN
	}
	return result
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
			Prop: NodeMetadataPropDimension,
		}, NodeMetadata{
			Prop: NodeMetadataPropMeasure,
		}, NodeMetadata{
			Prop: NodeMetadataPropAggregationFn,
		}, NodeMetadata{
			Prop: NodeMetadataPropDataType,
		}, NodeMetadata{
			Prop: NodeMetadataPropDescription,
		}, NodeMetadata{
			Prop: NodeMetadataPropDateFormat,
		})
	}
	for i := 0; i < len(metadata); i++ {
		metadata[i].DatasetID = n.DatasetID
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = string(c.Word)
		} else if metadata[i].Prop == NodeMetadataPropName {
			metadata[i].Value = c.Name
		} else if metadata[i].Prop == NodeMetadataPropDimension {
			if c.Dimension {
				metadata[i].Value = NodeMetadataPropValueTrue
			} else {
				metadata[i].Value = NodeMetadataPropValueFalse
			}
		} else if metadata[i].Prop == NodeMetadataPropMeasure {
			if c.Measure {
				metadata[i].Value = NodeMetadataPropValueTrue
			} else {
				metadata[i].Value = NodeMetadataPropValueFalse
			}
		} else if metadata[i].Prop == NodeMetadataPropAggregationFn {
			if _, ok := NodeMetadataAggregationFns[c.AggregationFn]; ok {
				metadata[i].Value = c.AggregationFn
			} else {
				metadata[i].Value = interpreter.AggregationFnCount
			}
		} else if metadata[i].Prop == NodeMetadataPropDataType {
			if _, ok := NodeMetadataDataTypes[c.DataType]; ok {
				metadata[i].Value = c.DataType
			} else {
				metadata[i].Value = interpreter.DataTypeString
			}
		} else if metadata[i].Prop == NodeMetadataPropDescription {
			metadata[i].Value = c.Description
		} else if metadata[i].Prop == NodeMetadataPropDateFormat {
			metadata[i].Value = c.DateFormat
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

//TableNode returns table node converted form of the node
func (n Node) TableNode() interpreter.TableNode {
	name := ""
	word := ""
	description := ""
	defauldDateFieldUID := ""
	datastoreID := 0
	for _, v := range n.NodeMetadatas {
		if v.Prop == NodeMetadataPropWord {
			word = v.Value
		} else if v.Prop == NodeMetadataPropName {
			name = v.Value
		} else if v.Prop == NodeMetadataPropDefaultDateFieldUID {
			defauldDateFieldUID = v.Value
		} else if v.Prop == NodeMetadataPropDescription {
			description = v.Value
		} else if v.Prop == NodeMetadataPropDatastoreID {
			datastoreID, _ = strconv.Atoi(v.Value)
		}
	}
	return interpreter.TableNode{
		UID:                 n.UID.String(),
		Word:                []rune(word),
		PUID:                n.PUID.String(),
		Name:                name,
		Children:            []interpreter.ColumnNode{},
		DefaultDateFieldUID: defauldDateFieldUID,
		Description:         description,
		DatastoreID:         uint(datastoreID),
	}
}

//FromTable converts the interpreter table node to node
func (n Node) FromTable(t interpreter.TableNode) Node {
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
			Prop: NodeMetadataPropDefaultDateFieldUID,
		}, NodeMetadata{
			Prop: NodeMetadataPropDescription,
		}, NodeMetadata{
			Prop: NodeMetadataPropDatastoreID,
		})
	}
	for i := 0; i < len(metadata); i++ {
		metadata[i].DatasetID = n.DatasetID
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = string(t.Word)
		} else if metadata[i].Prop == NodeMetadataPropName {
			metadata[i].Value = t.Name
		} else if metadata[i].Prop == NodeMetadataPropDescription {
			metadata[i].Value = t.Description
		} else if metadata[i].Prop == NodeMetadataPropDefaultDateFieldUID {
			metadata[i].Value = t.DefaultDateFieldUID
		} else if metadata[i].Prop == NodeMetadataPropDatastoreID {
			metadata[i].Value = strconv.Itoa(int(t.DatastoreID))
		}
	}
	uid, _ := uuid.Parse(t.UID)
	puid, _ := uuid.Parse(t.PUID)
	return Node{
		Model:         n.Model,
		UID:           uid,
		Type:          t.Type(),
		PUID:          puid,
		DatasetID:     n.DatasetID,
		NodeMetadatas: metadata,
	}
}

//KnowledgeBaseNode returns the knowledgebase node converted form of the node
func (n Node) KnowledgeBaseNode() interpreter.KnowledgeBaseNode {
	name := ""
	word := ""
	description := ""
	kbType := interpreter.SystemKB
	for _, v := range n.NodeMetadatas {
		if v.Prop == NodeMetadataPropWord {
			word = v.Value
		} else if v.Prop == NodeMetadataPropName {
			name = v.Value
		} else if v.Prop == NodeMetadataPropDescription {
			description = v.Value
		} else if v.Prop == NodeMetadataPropKBType && v.Value == NodeMetadataPropValueUserKB {
			kbType = interpreter.UserKB
		}
	}
	return interpreter.KnowledgeBaseNode{
		UID:         n.UID.String(),
		Word:        []rune(word),
		Name:        name,
		Children:    []interpreter.Node{},
		Description: description,
		KBType:      kbType,
	}
}

//FromKnowledgeBase converts the interpreter knowledgebase node to node
func (n Node) FromKnowledgeBase(k interpreter.KnowledgeBaseNode) Node {
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
			Prop: NodeMetadataPropDescription,
		}, NodeMetadata{
			Prop: NodeMetadataPropKBType,
		})
	}
	for i := 0; i < len(metadata); i++ {
		metadata[i].DatasetID = n.DatasetID
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = string(k.Word)
		} else if metadata[i].Prop == NodeMetadataPropName {
			metadata[i].Value = k.Name
		} else if metadata[i].Prop == NodeMetadataPropDescription {
			metadata[i].Value = k.Description
		} else if metadata[i].Prop == NodeMetadataPropKBType {
			metadata[i].Value = strconv.Itoa(int(k.KBType))
		}
	}
	uid, _ := uuid.Parse(k.UID)
	return Node{
		Model:         n.Model,
		UID:           uid,
		Type:          k.Type(),
		NodeMetadatas: metadata,
	}
}

//OperatorNode returns the operator node converted form of the node
func (n Node) OperatorNode() interpreter.OperatorNode {
	word := ""
	operation := ""
	for _, v := range n.NodeMetadatas {
		if v.Prop == NodeMetadataPropWord {
			word = v.Value
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueEqOperator {
			operation = interpreter.EqOperator
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueNotEqOperator {
			operation = interpreter.NotEqOperator
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueGreaterOperator {
			operation = interpreter.GreaterOperator
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueLessOperator {
			operation = interpreter.LessOperator
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueContainsOperator {
			operation = interpreter.ContainsOperator
		} else if v.Prop == NodeMetadataPropDimension && v.Value == NodeMetadataPropValueLikeOperator {
			operation = interpreter.LikeOperator
		}
	}
	result := interpreter.OperatorNode{
		UID:       n.UID.String(),
		Word:      []rune(word),
		PUID:      n.PUID.String(),
		Operation: operation,
	}
	if n.Parent != nil && n.PUID.String() == n.Parent.UID.String() {
		pN, ok := n.Parent.InterpreterNode()
		if ok {
			result.PN = pN
		}
	}
	return result
}

//FromOperatorNode converts the interpreter operator node to node
func (n Node) FromOperatorNode(o interpreter.OperatorNode) Node {
	metadata := []NodeMetadata{}
	for _, v := range n.NodeMetadatas {
		metadata = append(metadata, v)
	}
	if len(metadata) == 0 {
		metadata = append(metadata, NodeMetadata{
			Prop: NodeMetadataPropWord,
		}, NodeMetadata{
			Prop: NodeMetadataPropOperation,
		})
	}
	for i := 0; i < len(metadata); i++ {
		metadata[i].DatasetID = n.DatasetID
		if metadata[i].Prop == NodeMetadataPropWord {
			metadata[i].Value = string(o.Word)
		} else if metadata[i].Prop == NodeMetadataPropOperation {
			metadata[i].Value = o.Operation
		}
	}
	uid, _ := uuid.Parse(o.UID)
	puid, _ := uuid.Parse(o.PUID)
	return Node{
		Model:         n.Model,
		UID:           uid,
		Type:          o.Type(),
		PUID:          puid,
		NodeMetadatas: metadata,
	}
}
