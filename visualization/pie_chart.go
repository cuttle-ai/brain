// Copyright 2019 Cuttle.ai. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package visualization

import (
	"strings"

	"github.com/cuttle-ai/octopus/interpreter"
)

/*
 * This file contains the pie chart visualization definition
 */

//PieChart will return a pie visualization for a given query
func PieChart(q *interpreter.Query) Visualization {
	/*
	 * We will iterate through the select columns and add them to metrics and description
	 * We will iterate through the group by columns and add them to metrics and description
	 */
	var description strings.Builder
	result := Visualization{Title: "Pie Chart", Metrics: []Metric{}, Type: PieChartType}
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
