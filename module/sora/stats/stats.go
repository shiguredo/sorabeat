package stats

import (
	"encoding/json"
	"io/ioutil"
	"math"

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
	defaultScheme     = "http"
	httpPath          = "/"
	httpMethod        = "POST"
	targetHeaderKey   = "x-sora-target"
	targetHeaderValue = "Sora_20171010.GetStatsReport"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   httpPath,
	}.Build()

	float_array_keys = []string{"active_tasks", "active_tasks_all",
		"run_queue_lengths", "run_queue_lengths_all"}
)

// MetricSet type defines all fields of the MetricSet
// As a minimum it must inherit the mb.BaseMetricSet fields, but can be extended with
// additional entries. These variables can be used to persist data or configuration between
// multiple fetch calls.
type MetricSet struct {
	mb.BaseMetricSet
	http *helper.HTTP
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

	// erlang_vm フィールドの数値リストからいくつかフィールドを追加する
	if val, ok := stats["erlang_vm"]; ok {
		erlang_vm, _ := val.(map[string]interface{})
		statistics, _ := erlang_vm["statistics"].(map[string]interface{})
		for _, key := range float_array_keys {
			addStats(key, statistics)
		}
	}

	return stats, nil
}

func addStats(key string, m map[string]interface{}) {
	value, has_key := m[key]
	if !has_key {
		return
	}

	var numbers []float64
	for _, v := range value.([]interface{}) {
		numbers = append(numbers, v.(float64))
	}

	if len(numbers) == 0 {
		return
	}

	mean := mean(numbers)
	m[key+"_mean"] = mean
	m[key+"_stddev"] = calcStdDev(numbers, mean)
	min, max := MinMax(numbers)
	m[key+"_min"] = min
	m[key+"_max"] = max
	// TODO: 偏りの良い指標、ひとまず Max / Min でいく
	m[key+"_imbalance"] = max / math.Max(min, 1.)
}

func MinMax(numbers []float64) (min float64, max float64) {
	max = numbers[0]
	min = numbers[0]
	for _, value := range numbers {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func mean(numbers []float64) float64 {
	total := calcTotal(numbers)
	return total / float64(len(numbers))
}

func calcTotal(numbers []float64) (total float64) {
	for _, x := range numbers {
		total += x
	}
	return total
}

func calcStdDev(numbers []float64, mean float64) float64 {
	total := 0.0
	for _, number := range numbers {
		total += math.Pow(number-mean, 2)
	}
	variance := total / float64(len(numbers))
	return math.Sqrt(variance)
}
