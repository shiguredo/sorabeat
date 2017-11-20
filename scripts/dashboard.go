// Copyright 2017 Shiguredo Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
)

func main() {
	buf, readErr := readSoraFields("scripts/sora_fields.yml")
	if readErr != nil {
		debugPrint(readErr)
		os.Exit(1)
	}

	processErr := processRootNodes(buf)
	if processErr != nil {
		debugPrint(processErr)
		os.Exit(2)
	}

	debugPrint("SUCCEEDED!! ＼（＾ ＾）／")
}

func readSoraFields(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	return data, err
}

type RootNode struct {
	Key         string
	Title       string
	Type        string
	Description string
	Fields      []Node `yaml:"fields,omitempty"`
}

type Node struct {
	Name        string
	Title       string
	Type        string
	Cumulative  bool `yaml:"cumulative,omitempty"`
	Description string
	Fields      []Node `yaml:"fields,omitempty"`
}

func processRootNodes(buf []byte) error {
	var rootNodes []RootNode
	err1 := yaml.Unmarshal(buf, &rootNodes)
	if err1 != nil {
		return err1
	}
	visualizations := make([]map[string]interface{}, 0)
	for _, rootNode := range rootNodes {
		if rootNode.Key == "sora" {
			err := processSoraNode(rootNode, &visualizations)
			if err != nil {
				return err
			}
		}
	}
	vis1Json := jsonObj()
	vis1Json["objects"] = visualizations
	vis1Json["version"] = "1.0.0"
	vis1JsonBytes, marshalErr := json.Marshal(vis1Json)
	if marshalErr != nil {
		return marshalErr
	}
	print(string(vis1JsonBytes[:]))
	return nil
}

func processSoraNode(sora RootNode, visualizations *[]map[string]interface{}) error {
	for _, node := range sora.Fields {
		if node.Name == "connections" {
			err := processConnectionsNode(node, visualizations)
			if err != nil {
				return err
			}
		} else if node.Name == "stats" {
			err := processStatsNode(node, visualizations)
			if err != nil {
				return err
			}
		} else {
			// nop
		}
	}
	return nil
}

func processConnectionsNode(connections Node, visualizations *[]map[string]interface{}) error {
	for _, field := range connections.Fields {
		prefix := "sora.connections."
		splitMode := "terms"
		termsField := "sora.connections.channel_client_id"
		visualization, err := processField(prefix, field, splitMode, termsField)
		if err != nil {
			return err
		}
		// debugPrint(visualization)
		if visualization != nil {
			*visualizations = append(*visualizations, visualization)
		}
	}
	return nil
}

func processStatsNode(stats Node, visualizations *[]map[string]interface{}) error {
	for _, field := range stats.Fields {
		prefix := "sora.stats."
		splitMode := "everything"
		termsField := ""
		visualization, err := processField(prefix, field, splitMode, termsField)
		if err != nil {
			return err
		}
		// debugPrint(visualization)
		if visualization != nil {
			*visualizations = append(*visualizations, visualization)
		}
	}
	return nil
}

func processField(prefix string, field Node, splitMode string,
	termsField string) (map[string]interface{}, error) {

	if field.Type != "bytes" && field.Type != "long" {
		return nil, nil
	}

	item := field.Name
	metricsType := "max"
	var formatter string
	if field.Type == "bytes" {
		formatter = field.Type
	} else {
		formatter = "number"
	}
	axisMin := "0"
	derivative := field.Cumulative
	visualization, err := visualizationJson("sorabeat-vis1-", prefix, item,
		splitMode, derivative, termsField,
		metricsType, formatter, axisMin)
	return visualization, err
}

// visualization:
//     {
//       "id": "92381420-c525-11e7-b277-79c0643bd2c8",
//       "type": "visualization",
//       "version": 2,
//       "attributes": {
//         "title": "Sora/BEAM active_tasks_all",
//         "visState": "[下記参照:stringify された JSON]",
//         "uiStateJSON": "{}",
//         "description": "",
//         "version": 1,
//         "kibanaSavedObjectMeta": {
//           "searchSourceJSON": "{}"
//         }
//       }
//     },

// visState:
// {
//     "aggs": [],
//     "params": {
//         "axis_formatter": "number",
//         "axis_position": "left",
//         "id": "61ca57f0-469d-11e7-af02-69e470af7417",
//         "index_pattern": "*",
//         "interval": "auto",
//         "series": [
//             {
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "number",
//                 "id": "61ca57f1-469d-11e7-af02-69e470af7417",
//                 "label": "max",
//                 "line_width": 1,
//                 "metrics": [
//                     {
//                         "field": "sora.stats.erlang_vm.statistics.active_tasks_all_max",
//                         "id": "61ca57f2-469d-11e7-af02-69e470af7417",
//                         "type": "avg"
//                     }
//                 ],
//                 "point_size": 1,
//                 "seperate_axis": 0,
//                 "split_mode": "everything",
//                 "stacked": "none"
//             }
//         ],
//         "show_grid": 1,
//         "show_legend": 1,
//         "time_field": "@timestamp",
//         "type": "timeseries"
//     },
//     "title": "Sora/BEAM active_tasks_all",
//     "type": "metrics"
// }

// derivative 型の visState:
// {
//     "aggs": [],
//     "params": {
//         "axis_formatter": "number",
//         "axis_min": "0",
//         "axis_position": "left",
//         "id": "61ca57f0-469d-11e7-af02-69e470af7417",
//         "ignore_global_filter": 0,
//         "index_pattern": "*",
//         "interval": "auto",
//         "series": [
//             {
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "bytes",
//                 "id": "e8f96550-bfaf-11e7-ba99-7dd83649120a",
//                 "label": "sent",
//                 "line_width": "2",
//                 "metrics": [
//                     {
//                         "field": "sora.connections.rtp.total_sent_bytes",
//                         "id": "e8f96551-bfaf-11e7-ba99-7dd83649120a",
//                         "type": "max"
//                     },
//                     {
//                         "field": "e8f96551-bfaf-11e7-ba99-7dd83649120a",
//                         "id": "f6cf9230-bfaf-11e7-ba99-7dd83649120a",
//                         "type": "derivative",
//                         "unit": "1s"
//                     }
//                 ],
//                 "point_size": "2",
//                 "seperate_axis": 0,
//                 "split_mode": "terms",
//                 "stacked": "none",
//                 "terms_field": "sora.connections.channel_client_id",
//                 "terms_order_by": "e8f96551-bfaf-11e7-ba99-7dd83649120a",
//                 "terms_size": "10",
//                 "value_template": "{{value}}/s"
//             }
//         ],
//         "show_grid": 1,
//         "show_legend": 1,
//         "time_field": "@timestamp",
//         "type": "timeseries"
//     },
//     "title": "Sora total bytes",
//     "type": "metrics"
// }

// derivative 型、かつ sum aggregation の visState (TODO: 未だ作ってない)
// {
//     "aggs": [],
//     "params": {
//         "axis_formatter": "number",
//         "axis_min": "0",
//         "axis_position": "left",
//         "id": "61ca57f0-469d-11e7-af02-69e470af7417",
//         "index_pattern": "*",
//         "interval": "auto",
//         "series": [
//             {
//                 "axis_min": "0",
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "bytes",
//                 "id": "a9b0f6f0-c370-11e7-9e32-ff5b8223c99f",
//                 "label": "sent (sum)",
//                 "line_width": "3",
//                 "metrics": [
//                     {
//                         "field": "sora.connections.rtp.total_sent_bytes",
//                         "id": "a9b0f6f1-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "max"
//                     },
//                     {
//                         "field": "a9b0f6f1-c370-11e7-9e32-ff5b8223c99f",
//                         "id": "a9b0f6f2-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "derivative",
//                         "unit": "1s"
//                     },
//                     {
//                         "function": "sum",
//                         "id": "a9b0f6f3-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "series_agg"
//                     }
//                 ],
//                 "point_size": "4",
//                 "seperate_axis": 0,
//                 "split_filters": [
//                     {
//                         "color": "#68BC00",
//                         "id": "074842f0-c36c-11e7-9cc7-5705c84b2ed3"
//                     }
//                 ],
//                 "split_mode": "terms",
//                 "stacked": "none",
//                 "terms_field": "sora.connections.channel_client_id",
//                 "terms_order_by": "a9b0f6f1-c370-11e7-9e32-ff5b8223c99f",
//                 "terms_size": "1000",
//                 "value_template": "{{value}}/s"
//             },
//             {
//                 "axis_min": "0",
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "bytes",
//                 "id": "55703560-c370-11e7-9e32-ff5b8223c99f",
//                 "label": "recieved (sum)",
//                 "line_width": "3",
//                 "metrics": [
//                     {
//                         "field": "sora.connections.rtp.total_received_bytes",
//                         "id": "55703561-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "max"
//                     },
//                     {
//                         "field": "55703561-c370-11e7-9e32-ff5b8223c99f",
//                         "id": "55705c70-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "derivative",
//                         "unit": "1s"
//                     },
//                     {
//                         "function": "sum",
//                         "id": "5bb8f150-c370-11e7-9e32-ff5b8223c99f",
//                         "type": "series_agg"
//                     }
//                 ],
//                 "point_size": "4",
//                 "seperate_axis": 0,
//                 "split_filters": [
//                     {
//                         "color": "#68BC00",
//                         "id": "074842f0-c36c-11e7-9cc7-5705c84b2ed3"
//                     }
//                 ],
//                 "split_mode": "terms",
//                 "stacked": "none",
//                 "terms_field": "sora.connections.channel_client_id",
//                 "terms_order_by": "55703561-c370-11e7-9e32-ff5b8223c99f",
//                 "terms_size": "1000",
//                 "value_template": "{{value}}/s"
//             }
//         ],
//         "show_grid": 1,
//         "show_legend": 1,
//         "time_field": "@timestamp",
//         "type": "timeseries"
//     },
//     "title": "Sora total bytes (aggregated/sum)",
//     "type": "metrics"
// }

func visualizationJson(
	idPrefix string,
	prefix string, item string,
	splitMode string, derivative bool, termsField string,
	metricsType string,
	formatter string, axisMin string) (map[string]interface{}, error) {
	title := item + " [Sorabeat]"
	values := jsonObj()
	{
		values["id"] = idPrefix + prefix + item
		values["type"] = "visualization"
		values["version"] = 1
		{
			attrs := jsonObj()
			attrs["title"] = title
			{
				visState := jsonObj()
				visState["aggs"] = make([]map[string]interface{}, 0)
				{
					params := jsonObj()
					params["axis_formatter"] = "number"
					params["axis_min"] = axisMin
					params["axis_position"] = "left"
					params["id"] = uuid.NewV4()
					params["index_pattern"] = "*"
					params["interval"] = "auto"

					{
						series := make([]map[string]interface{}, 0)
						series0 := jsonObj()

						series0["axis_position"] = "right"
						series0["chart_type"] = "line"
						series0["color"] = "#68BC00"
						series0["fill"] = "0"
						series0["formatter"] = formatter // e.g. "number", "bytes"
						series0["id"] = uuid.NewV4()
						series0["label"] = item
						series0["line_width"] = "2"

						metrics0Id := uuid.NewV4()
						{
							metrics := make([]map[string]interface{}, 0)
							metrics0 := jsonObj()
							metrics0["field"] = prefix + item
							metrics0["id"] = metrics0Id
							metrics0["type"] = metricsType // e.g. "avg", "max"
							metrics = append(metrics, metrics0)
							if derivative {
								metrics1 := jsonObj()
								metrics1Id := uuid.NewV4()
								metrics1["id"] = metrics1Id
								metrics1["type"] = "derivative"
								metrics1["field"] = metrics0Id
								metrics1["unit"] = "1s"
								metrics = append(metrics, metrics1)
							}
							series0["metrics"] = metrics
						}

						series0["point_size"] = "2"
						series0["seperate_axis"] = 0
						series0["stacked"] = "none"

						series0["split_mode"] = splitMode // "everything" or "terms"
						if splitMode == "terms" {
							series0["terms_field"] = termsField // e.g. "sora.connections.channel_client_id"
							series0["terms_order_by"] = metrics0Id
							series0["terms_size"] = 20 // TODO: sufficient??
							series0["value_template"] = "{{value}}/s"
						}

						series = append(series, series0)
						params["series"] = series
					}

					params["show_grid"] = 1
					params["show_legend"] = 1
					params["time_field"] = "@timestamp"
					params["type"] = "timeseries"
					visState["params"] = params
				}
				visState["title"] = title
				visState["type"] = "metrics"

				visStateBytes, _ := json.Marshal(visState)
				attrs["visState"] = string(visStateBytes[:])
			}

			attrs["uiStateJSON"] = "{}"
			attrs["description"] = ""
			attrs["version"] = 1
			{
				kibanaSavedObjectMeta := jsonObj()
				kibanaSavedObjectMeta["searchSourceJSON"] = "{}"
				attrs["kibanaSavedObjectMeta"] = kibanaSavedObjectMeta
			}
			values["attributes"] = attrs
		}
	}
	return values, nil
}

func jsonObj() map[string]interface{} {
	return make(map[string]interface{})
}

func print(arg interface{}) {
	fmt.Printf("%s\n", arg)
}

func debugPrintf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args)
}

func debugPrint(arg interface{}) {
	fmt.Fprintf(os.Stderr, "%#v\n", arg)
}
