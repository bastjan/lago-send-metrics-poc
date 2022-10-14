package main

import (
	"github.com/appuio/appuio-cloud-reporting/pkg/sourcekey"
	"github.com/prometheus/common/model"

	"github.com/bastjan/lago-send-metrics-poc/prommock"
)

var promData = prommock.FakeQuerier{
	Queries: map[string]map[string]model.SampleValue{
		"PRODUCTS": {
			sourcekey.SourceKey{
				Query:     "tantanmen",
				Zone:      "langstrasse",
				Tenant:    "horny-for-ramen",
				Namespace: "takeout",
				Class:     "vegan",
			}.String(): 500,
		},
	},
}
