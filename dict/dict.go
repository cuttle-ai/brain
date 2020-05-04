// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package dict has the implementation of the dictionary api for the platform
package dict

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cuttle-ai/brain/log"
	"github.com/cuttle-ai/brain/models"
	"github.com/cuttle-ai/octopus/interpreter"
	"github.com/jinzhu/gorm"
)

//DAgg is the dict aggregator for getting the dict from the database
type DAgg struct {
	db *gorm.DB
	l  log.Log
}

//NewDAgg returns an instance of DAgg dict aggregator
func NewDAgg(db *gorm.DB, l log.Log) *DAgg {
	return &DAgg{db, l}
}

//Get returns the user dictionary from the database
func (d DAgg) Get(ID string, update bool) (interpreter.DICT, error) {
	/*
	 * We will convert the id to integer
	 * We will get all the datasets the user has access to
	 * Then we will get the nodes belonging to those datasets
	 * Then we will add the system dict
	 */
	result := interpreter.DICT{Map: map[string]interpreter.Token{}}
	//parsing the user id
	id, err := strconv.Atoi(ID)
	if err != nil {
		d.l.Error("error while parsing the id from string to integer", ID)
		return result, err
	}

	//finding the datasets the user has access to
	datasets := []models.DatsetUserMapping{}
	err = d.db.Where("user_id = ?", id).Find(&datasets).Error
	if err != nil {
		d.l.Error("error while getting the list of datasets the user has access to", ID)
		return result, err
	}

	//iterating through the list and getting the datasets
	for _, v := range datasets {
		req := DatasetRequest{ID: strconv.Itoa(int(v.DatasetID)), SubscribeID: ID, Type: DatasetGet, Out: make(chan DatasetRequest)}
		if update {
			req.Type = DatasetUpdate
		}
		DatasetInputChannel <- req
		req = <-req.Out
		if !req.Valid {
			continue
		}
		//iterating through the result and adding to the token list
		for k, t := range req.Dataset.D {
			existing, ok := result.Map[k]
			if ok {
				t.Nodes = append(t.Nodes, existing.Nodes...)
			}
			result.Map[k] = t
		}
	}

	//adding the system dict
	systemnDict := SystemDICT()
	for k, v := range systemnDict.Map {
		existing, ok := result.Map[k]
		if ok {
			v.Nodes = append(v.Nodes, existing.Nodes...)
		}
		result.Map[strings.ToLower(k)] = v
	}

	return result, nil
}

//GetDataset will get the dataset required for the given id
func (d DAgg) GetDataset(ID string) (Dataset, error) {
	/*
	 * We will parse the id of the dataset
	 * We will find all the nodes associated with the dataset
	 * We will find all the node metadata associated with the dataset
	 * Will convert them into token
	 */
	result := Dataset{D: map[string]interpreter.Token{}}
	//parsing the id of the dataset
	id, err := strconv.Atoi(ID)
	if err != nil {
		d.l.Error("error while parsing the id from string to integer", ID)
		return result, err
	}

	//finding all the nodes associated with the dataset
	nodes := []models.Node{}
	err = d.db.Where("dataset_id = ?", id).Find(&nodes).Error
	if err != nil {
		d.l.Error("error while getting the list of nodes the dataset has access to", ID)
		return result, err
	}

	//finding all the node metadata associated with the dataset
	nodeMetadatas := []models.NodeMetadata{}
	err = d.db.Where("dataset_id = ?", id).Find(&nodeMetadatas).Error
	if err != nil {
		d.l.Error("error while getting the list of node metadata the dataset has access to", ID)
		return result, err
	}

	//converting the nodes to tokens
	//we will iterate through the nodes and store them in map
	//then we will iterate through the node metadatas and store them to the correspoding nodes in the map
	//finally will convert them to interpreter nodes
	nMap := map[uint]models.Node{}
	var tableNode *models.Node
	for _, n := range nodes {
		nMap[n.ID] = n
	}
	for _, m := range nodeMetadatas {
		n, ok := nMap[m.NodeID]
		if !ok {
			continue
		}
		if len(n.NodeMetadatas) == 0 {
			n.NodeMetadatas = []models.NodeMetadata{}
		}
		n.NodeMetadatas = append(n.NodeMetadatas, m)
		nMap[m.NodeID] = n
		if n.Type == interpreter.Table {
			tableNode = &n
		}
	}
	//if default date exists add it to the node
	if tableNode != nil && len(tableNode.TableNode().DefaultDateFieldUID) > 0 {
		tConverted := tableNode.TableNode()
		for _, n := range nMap {
			if n.Type != interpreter.Table && tConverted.DefaultDateFieldUID == n.UID.String() {
				c := n.ColumnNode()
				tConverted.DefaultDateField = &c
				*tableNode = tableNode.FromTable(tConverted)
				break
			}
		}
	}
	for _, n := range nMap {
		if n.Type != interpreter.Table && tableNode != nil {
			n.Parent = tableNode
			n.PUID = tableNode.UID
		}
		iN, ok := n.InterpreterNode()
		if !ok {
			continue
		}
		tok, ok := result.D[string(iN.TokenWord())]
		if !ok {
			tok = interpreter.Token{Word: iN.TokenWord(), Nodes: []interpreter.Node{}}
		}
		tok.Nodes = append(tok.Nodes, iN)
		result.D[strings.ToLower(string(iN.TokenWord()))] = tok
	}

	return result, nil
}

//DatasetRequestType is the type of the request for the dataset
type DatasetRequestType uint

const (
	//DatasetUpdate comes to update and datatset and update corresponding subscribed ids
	DatasetUpdate DatasetRequestType = 1
	//DatasetGet returns the dataset of a given id
	DatasetGet DatasetRequestType = 2
	//DatasetRemove to remove the dataset from the cache
	DatasetRemove DatasetRequestType = 3
)

//DatasetClearCheckInterval is the interval after which the datatset removal check has to run
const DatasetClearCheckInterval = time.Minute * 20

//DatasetExpiry is the expiry time after which the datatset expiries without any active usage
const DatasetExpiry = time.Hour * 4

//DatasetAggregator is the aggregator to get the dict from a service or database
type DatasetAggregator interface {
	GetDataset(ID string) (Dataset, error)
}

//defaultAggregator to be used in the dictionary
var defaultAggregator aggregator

type aggregator struct {
	agg DatasetAggregator
	m   sync.Mutex
}

//SetDefaultDatasetAggregator sets the default aggregator as the passed param
func SetDefaultDatasetAggregator(agg DatasetAggregator) {
	defaultAggregator.m.Lock()
	defaultAggregator.agg = agg
	defaultAggregator.m.Unlock()
}

func getDataset(ID string) (Dataset, bool) {
	defaultAggregator.m.Lock()
	if defaultAggregator.agg == nil {
		return Dataset{}, false
	}
	d, err := defaultAggregator.agg.GetDataset(ID)
	if err != nil {
		return Dataset{}, false
	}
	defaultAggregator.m.Unlock()
	return d, true
}

//Dataset is the dataset instance having the node tokens
type Dataset struct {
	D        map[string]interpreter.Token
	LastUsed time.Time
}

//DatasetRequest can be used to make a request to get the dataset cache
type DatasetRequest struct {
	//ID of the dataset
	ID string
	//SubscribeID is the id which subscribes to the dataset
	SubscribeID string
	//Type is the type of the dictionary request. It can have Add, Get, Remove
	Type DatasetRequestType
	//Dataset has the tokens mapped to the token string for the dataset
	Dataset Dataset
	//Valid indicates that the dict is valid. During get requests, if valid is false then cache couldn't find the dict
	Valid bool
	//Out channel for sending response to the requester
	Out chan DatasetRequest
}

//DatasetInputChannel is the input channel to communicate with the cache
var DatasetInputChannel chan DatasetRequest

//SendDatasetToChannel sends a dataset request to the channel. This function is to be used with go routines so that
//datasets isn't blocked by the requests
func SendDatasetToChannel(ch chan DatasetRequest, req DatasetRequest) {
	ch <- req
}

func init() {
	DatasetInputChannel = make(chan DatasetRequest)
	defaultAggregator = aggregator{}
	go Datasets(DatasetInputChannel)
	go cacheClearCheck(DatasetInputChannel)
}

//Datasets is the cache providing the datatsets. When a dataset is updated, coresponding users who all have
//access to that dataset get their DICTs updated automatically
func Datasets(in chan DatasetRequest) {
	subscribedMap := map[string][]string{}
	datasets := map[string]Dataset{}
	for {
		req := <-in
		switch req.Type {
		case DatasetGet:
			req.Dataset, req.Valid = datasets[req.ID]
			if !req.Valid {
				req.Dataset, req.Valid = getDataset(req.ID)
			}
			req.Dataset.LastUsed = time.Now()
			datasets[req.ID] = req.Dataset
			s, ok := subscribedMap[req.ID]
			if !ok {
				s = []string{}
			}
			s = append(s, req.SubscribeID)
			subscribedMap[req.ID] = s
			go SendDatasetToChannel(req.Out, req)
			break
		case DatasetUpdate:
			delete(datasets, req.ID)
			req.Dataset, req.Valid = getDataset(req.ID)
			req.Dataset.LastUsed = time.Now()
			v, _ := subscribedMap[req.ID]
			for _, k := range v {
				go interpreter.SendDICTToChannel(interpreter.DICTInputChannel, interpreter.DICTRequest{ID: k, Type: interpreter.DICTRemove})
			}
			go SendDatasetToChannel(req.Out, req)
			break
		case DatasetRemove:
			//we will iterate over the cache and check the last usage
			t := time.Now()
			for k, v := range datasets {
				if v.LastUsed.Add(DatasetExpiry).After(t) {
					delete(datasets, k)
					delete(subscribedMap, k)
				}
			}
		}
	}
}

func cacheClearCheck(in chan DatasetRequest) {
	for {
		time.Sleep(DatasetClearCheckInterval)
		go SendDatasetToChannel(in, DatasetRequest{Type: DatasetRemove})
	}
}

//SystemDICT returns the system dictionary available for all the users
func SystemDICT() interpreter.DICT {
	d := map[string]interpreter.Token{
		"is": {
			Word:  []rune("is"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "equal-is", Word: []rune("is"), Operation: interpreter.EqOperator}},
		},
		"not": {
			Word:  []rune("not"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "not-equal", Word: []rune("not"), Operation: interpreter.NotEqOperator}},
		},
		"<": {
			Word:  []rune("<"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "less-than", Word: []rune("<="), Operation: interpreter.LessOperator}},
		},
		">": {
			Word:  []rune(">"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "greater-than", Word: []rune(">="), Operation: interpreter.GreaterOperator}},
		},
		"less than": {
			Word:  []rune("less than"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "less-than", Word: []rune("<="), Operation: interpreter.LessOperator}},
		},
		"greater than": {
			Word:  []rune("greater than"),
			Nodes: []interpreter.Node{&interpreter.OperatorNode{UID: "greater-than", Word: []rune(">="), Operation: interpreter.GreaterOperator}},
		},
	}
	return interpreter.DICT{Map: d}
}
