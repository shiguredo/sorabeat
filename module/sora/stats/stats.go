package stats

import (
	"encoding/json"
	"io/ioutil"

	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/metricbeat/helper"
	"github.com/elastic/beats/metricbeat/mb"
	"github.com/elastic/beats/metricbeat/mb/parse"
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	if err := mb.Registry.AddMetricSet("sora", "stats", New, hostParser); err != nil {
		panic(err)
	}
}

const (
	defaultScheme = "http"
	httpPath = "/"
	httpMethod = "POST"
	targetHeaderKey = "x-sora-target"
	targetHeaderValue = "Sora_20171010.GetStatsReport"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   httpPath,
	}.Build()
)

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	http         *helper.HTTP
}

// New create a new instance of the MetricSet
// Part of new is also setting up the configuration by processing additional
// configuration entries if needed.
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {

	config := struct{}{}

	if err := base.Module().UnpackConfig(&config); err != nil {
		return nil, err
	}

	http := helper.NewHTTP(base)
	http.SetMethod(httpMethod)
	http.SetHeader(targetHeaderKey, targetHeaderValue)
	return &MetricSet{
		BaseMetricSet: base,
		http:          http,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() (common.MapStr, error) {

	response, err := m.http.FetchResponse()
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	err = json.Unmarshal(body, &stats)
	if err != nil {
		return nil, err
	}

	// TODO: erlang run queue にフィールドを追加する
	return stats, nil
}
