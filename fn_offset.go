package main

import (
	"fmt"
	"math"

	"github.com/gogo/protobuf/proto"
)

// offset(seriesList,factor)
func offset(e *expr, from, until int32, values map[metricRequest][]*metricData) []*metricData {
	arg, err := getSeriesArg(e.args[0], from, until, values)
	if err != nil {
		return nil
	}
	factor, err := getFloatArg(e, 1)
	if err != nil {
		return nil
	}
	var results []*metricData

	for _, a := range arg {
		r := *a
		r.Name = proto.String(fmt.Sprintf("offset(%s,%g)", a.GetName(), factor))
		r.Values = make([]float64, len(a.Values))
		r.IsAbsent = make([]bool, len(a.Values))

		for i, v := range a.Values {
			if a.IsAbsent[i] {
				r.Values[i] = 0
				r.IsAbsent[i] = true
				continue
			}
			r.Values[i] = v + factor
		}
		results = append(results, &r)
	}
	return results
}

// offsetToZero(seriesList)
func offsetToZero(e *expr, from, until int32, values map[metricRequest][]*metricData) []*metricData {
	return forEachSeriesDo(e, from, until, values, func(a *metricData, r *metricData) *metricData {
		minimum := math.Inf(1)
		for i, v := range a.Values {
			if !a.IsAbsent[i] && v < minimum {
				minimum = v
			}
		}
		for i, v := range a.Values {
			if a.IsAbsent[i] {
				r.Values[i] = 0
				r.IsAbsent[i] = true
				continue
			}
			r.Values[i] = v - minimum
		}
		return r
	})
}
