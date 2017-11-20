// Copyright [yyyy] [name of copyright owner]
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
	buf, readErr := readSoraFields("util/sora_fields.yml")
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
	soraJson := jsonObj()
	soraJson["objects"] = visualizations
	soraJson["version"] = "1.0.2"
	jsonBytes, marshalErr := json.Marshal(soraJson)
	if marshalErr != nil {
		return marshalErr
	}
	print(string(jsonBytes[:]))
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
		if field.Type != "byte" && field.Type != "long" {
			continue
		}
		prefix := "sora.connections"
		item := field.Name
		metricsType := "max"
		formatter := field.Type
		axis_min := int32(0)
		splitMode := "terms"
		derivative := true
		termsField := "sora.connections.channel_client_id"
		visualization, err := visualizationJson(prefix, item,
			splitMode, derivative, termsField,
			metricsType, formatter, axis_min)
		if err != nil {
			return err
		}
		debugPrint(visualization)
		*visualizations = append(*visualizations, visualization)
	}
	return nil
}

func processStatsNode(Node, *[]map[string]interface{}) error {
	// TODO: NYI
	return nil
}

// visualization:
//     {
//       "id": "d0ec26d0-bea8-11e7-b277-79c0643bd2c8-3",
//       "type": "visualization",
//       "version": 1,
//       "attributes": {
//         "title": "3BEAM memory",
//         "visState": "[下記参照:stringify された JSON]",
//         "uiStateJSON": "{}",
//         "description": "",
//         "version": 1,
//         "kibanaSavedObjectMeta": {
//           "searchSourceJSON": "{}"
//         }
//       }

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
//                 "formatter": "bytes",
//                 "id": "61ca57f1-469d-11e7-af02-69e470af7417",
//                 "label": "beam_mem_total",
//                 "line_width": "2",
//                 "metrics": [
//                     {
//                         "field": "sora.stats.erlang_vm.memory.total",
//                         "id": "61ca57f2-469d-11e7-af02-69e470af7417",
//                         "type": "max"
//                     },
//                     {
//                         "function": "sum",
//                         "id": "5a9be470-c524-11e7-90ad-15a4935f7944",
//                         "type": "series_agg"
//                     }
//                 ],
//                 "point_size": "2",
//                 "seperate_axis": 0,
//                 "split_mode": "terms",
//                 "stacked": "none",
//                 "terms_field": "beat.hostname",
//                 "terms_size": "100"
//             },
//             {
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "bytes",
//                 "id": "b5f7f980-bea8-11e7-a725-b1c1e3e1f448",
//                 "label": "beam_mem_binary",
//                 "line_width": "2",
//                 "metrics": [
//                     {
//                         "field": "sora.stats.erlang_vm.memory.binary",
//                         "id": "b5f7f981-bea8-11e7-a725-b1c1e3e1f448",
//                         "type": "max"
//                     },
//                     {
//                         "function": "sum",
//                         "id": "763e0500-c524-11e7-90ad-15a4935f7944",
//                         "type": "series_agg"
//                     }
//                 ],
//                 "point_size": "2",
//                 "seperate_axis": 0,
//                 "split_mode": "terms",
//                 "stacked": "none",
//                 "terms_field": "beat.hostname",
//                 "terms_size": "100"
//             }
//         ],
//         "show_grid": 1,
//         "show_legend": 1,
//         "time_field": "@timestamp",
//         "type": "timeseries"
//     },
//     "title": "BEAM memory",
//     "type": "metrics"
// }

// derivative 型の visState:
// {
//     "aggs": [],
//     "params": {
//         "axis_formatter": "number",
//         "axis_min": "0",
//         "axis_position": "left",
//         "background_color_rules": [
//             {
//                 "id": "00e5cb30-c371-11e7-9e32-ff5b8223c99f"
//             }
//         ],
//         "bar_color_rules": [
//             {
//                 "id": "02da2120-c371-11e7-9e32-ff5b8223c99f"
//             }
//         ],
//         "gauge_color_rules": [
//             {
//                 "id": "065676a0-c371-11e7-9e32-ff5b8223c99f"
//             }
//         ],
//         "gauge_inner_width": 10,
//         "gauge_style": "half",
//         "gauge_width": 10,
//         "id": "61ca57f0-469d-11e7-af02-69e470af7417",
//         "ignore_global_filter": 0,
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
	prefix string, item string,
	splitMode string, derivative bool, termsField string,
	metricsType string,
	formatter string, axis_min int32) (map[string]interface{}, error) {
	title := "[Sora] " + item
	values := jsonObj()
	{
		values["_id"] = uuid.NewV4()
		values["_type"] = "visualization"
		{
			source := jsonObj()
			source["title"] = title
			{
				visState := jsonObj()
				visState["title"] = title
				visState["type"] = "metrics"
				{
					params := jsonObj()
					params["id"] = uuid.NewV4()
					params["type"] = "timeseries"
					series := make([]map[string]interface{}, 0)
					series0 := jsonObj()
					{
						series0["id"] = uuid.NewV4()
						series0["color"] = "#68BC00"
						// TODO: term で split する場合
						metrics := make([]map[string]interface{}, 0)
						metrics0Id := uuid.NewV4()
						{
							metrics0 := jsonObj()
							metrics0["id"] = metrics0Id
							metrics0["type"] = metricsType // e.g. "max"
							metrics0["field"] = prefix + item
							metrics = append(metrics, metrics0)
						}
						if derivative {
							metrics1 := jsonObj()
							metrics1["id"] = uuid.NewV4()
							metrics1["type"] = "derivative"
							metrics1["field"] = metrics0Id
							metrics1["unit"] = "1s"
							metrics = append(metrics, metrics1)
						}

						series0["split_mode"] = splitMode // "everything" or "terms"
						if splitMode == "terms" {
							series0["terms_field"] = termsField
							series0["terms_order_by"] = metrics0Id
							series0["terms_size"] = 20 // TODO: sufficient??
							series0["value_template"] = "{{value}}/s"
						}
						series0["mrtrics"] = metrics
						series0["seperate_axis"] = 0
						series0["axis_position"] = "left"
						series0["formatter"] = formatter
						series0["chart_type"] = "line"
						series0["line_width"] = 1
						series0["point_size"] = 1
						series0["fill"] = 0
						series0["stacked"] = "none"
						series0["label"] = item
					}
					series = append(series, series0)
					params["series"] = series
					visState["params"] = params
				}
				visState["time_field"] = "@timestamp"
				visState["index_pattern"] = "*"
				visState["interval"] = "auto"
				visState["axis_position"] = "left"
				visState["axis_formatter"] = "number"
				visState["show_legend"] = 1
				visState["show_grid"] = 1
				visState["axis_min"] = axis_min
				visState["aggs"] = make([]map[string]interface{}, 0)

				visStateBytes, _ := json.Marshal(visState)
				source["visState"] = string(visStateBytes[:])
			}
			values["_source"] = source
		}
		values["uiStateJson"] = "{}"
		values["description"] = title
		values["version"] = 1
		kibanaSavedObjectMeta := jsonObj()
		{
			kibanaSavedObjectMeta["searchSourceJSON"] = jsonObj()
		}
		values["kibanaSavedObjectMeta"] = kibanaSavedObjectMeta
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
