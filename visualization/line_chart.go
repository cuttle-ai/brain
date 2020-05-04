// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package visualization

import (
	"strings"

	"github.com/cuttle-ai/octopus/interpreter"
)

/*
 * This file contains the line chart visualization definition
 */

//LineChart will return a line visualization for a given query
func LineChart(q *interpreter.Query) Visualization {
	/*
	 * We will iterate through the select columns and add them to metrics and description
	 * We will iterate through the group by columns and add them to metrics and description
	 */
	var description strings.Builder
	result := Visualization{Title: "Line Chart", Metrics: []Metric{}, Type: LineChartType}
	//iterating through the select columns
	for k, v := range q.Select {
		//adding to metrics
		result.Metrics = append(result.Metrics, Metric{
			DisplayName: string(v.Word),
			Measure:     v.Measure,
			Dimension:   v.Dimension,
			Name:        v.Name,
			ResourceID:  v.UID,
		})
		//adding to the description
		if k != 0 {
			description.WriteString(", ")
		}
		description.WriteString(string(v.Word))
	}

	//iterating through the group by columns
	if len(q.GroupBy) > 0 {
		description.WriteString(" across ")
	}
	for k, v := range q.GroupBy {
		//adding to metrics
		result.Metrics = append(result.Metrics, Metric{
			DisplayName: string(v.Word),
			Measure:     v.Measure,
			Dimension:   v.Dimension,
			Name:        v.Name,
			ResourceID:  v.UID,
		})
		//adding to the description
		if k != 0 {
			description.WriteString(", ")
		}
		description.WriteString(string(v.Word))
	}

	result.Description = description.String()
	result.Title = result.Description
	return result
}
