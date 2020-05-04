// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

//Package visualization has the collection of visualizations and its utilities for the platform
package visualization

import "github.com/cuttle-ai/octopus/interpreter"

//Metric holds the information about an metric to be used in the visualization
type Metric struct {
	//ResourceID of the item which is represented by the metric
	ResourceID string `json:"resource_id,omitempty"`
	//DisplayName is the name to be used as a display
	DisplayName string `json:"display_name,omitempty"`
	//Name of the actual item in the data
	Name string `json:"name,omitempty"`
	//Measure flag states whether the metric is a measure type value
	Measure bool `json:"measure,omitempty"`
	//Dimension flag states whether the metric is a dimension type value
	Dimension bool `json:"dimension,omitempty"`
}

//Visualization has the information about a visualization
type Visualization struct {
	//Metrics holds the information about the metrics to be used in a visualization
	Metrics []Metric `json:"metrics,omitempty"`
	//Type indicates the type of the visualization
	Type string `json:"type,omitempty"`
	//Title of the visualization
	Title string `json:"title,omitempty"`
	//Description of the visualization
	Description string `json:"description,omitempty"`
}

const (
	//TableType is the table type of visualization
	TableType = "TABLE"
	//LineChartType is the line chart type of visualization
	LineChartType = "LINECHART"
)

//SuggestVisualization suggests the visualization to be used for the query
func SuggestVisualization(q *interpreter.Query) Visualization {
	/*
	 * If the no of select columns = 1 and group by columns = 1 we select line chart
	 * Default is table
	 */
	if len(q.Select) == 1 && len(q.GroupBy) == 1 {
		return LineChart(q)
	}
	return Table(q)
}
