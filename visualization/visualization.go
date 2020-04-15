// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package visualization has the collection of visualizations and its utilities for the platform
package visualization

import "github.com/cuttle-ai/octopus/interpreter"

//Metric holds the information about an metric to be used in the visualization
type Metric struct {
	//ResourceID of the item which is represented by the metric
	ResourceID string
	//DisplayName is the name to be used as a display
	DisplayName string
	//Name of the actual item in the data
	Name string
	//Measure flag states whether the metric is a measure type value
	Measure bool
	//Dimension flag states whether the metric is a dimension type value
	Dimension bool
}

//Visualization has the information about a visualization
type Visualization struct {
	//Metrics holds the information about the metrics to be used in a visualization
	Metrics []Metric
	//Type indicates the type of the visualization
	Type string
	//Title of the visualization
	Title string
	//Description of the visualization
	Description string
}

const (
	//TableType is the table type of visualization
	TableType = "TABLE"
)

//SuggestVisualization suggests the visualization to be used for the query
func SuggestVisualization(q *interpreter.Query) Visualization {
	//For now we are blindly suggesting table visualization
	return Table(q)
}
