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
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
	"github.com/satori/go.uuid"
)

func main() {
	buf, err1 := readSoraFields("util/sora_fields.yml")
	if err1 != nil {
		debugPrint(err1)
		os.Exit(1)
	}

	res, err2 := processRootNodes(buf)
	if err2 != nil {
		debugPrint(err2)
		os.Exit(2)
	}
	debugPrint(res)

	prefix := "sora.connections"
	item := "rtp.total_received_bytes"
	metricsType := "max"
	formatter := "byte"
	axis_min := int32(0)
	{
		splitMode := "everything"
		derivative := false
		termsField := ""
		vres, reserr := visualizationJson(prefix, item,
			splitMode, derivative, termsField,
			metricsType, formatter, axis_min)
		if reserr != nil {
			debugPrint(reserr)
			os.Exit(3)
		}
		print(vres)
	}

	{
		splitMode := "terms"
		derivative := true
		termsField := "sora.connections.channel_client_id"
		vres, reserr := visualizationJson(prefix, item,
			splitMode, derivative, termsField,
			metricsType, formatter, axis_min)
		if reserr != nil {
			debugPrint(reserr)
			os.Exit(3)
		}
		print(vres)
	}

	debugPrint("SUCCEEDED!! ＼（＾ ＾）／")
}

func readSoraFields(filePath string) ([]byte, error) {
    data, err := ioutil.ReadFile(filePath)
	return data, err
}

type RootNode struct {
	Key string
	Title string
	Type string
	Description string
	Fields []Node `yaml:"fields,omitempty"`
}

type Node struct {
	Name string
	Title string
	Type string
	Description string
	Fields []Node `yaml:"fields,omitempty"`
}

func processRootNodes(buf []byte) (interface{}, error) {
	var rootNodes []RootNode
	err1 := yaml.Unmarshal(buf, &rootNodes)
	if err1 != nil {
		return nil, err1
	}
	for _, rootNode := range rootNodes {
		if rootNode.Key == "sora" {
			processSoraNode(rootNode)
		}
	}

	if 1 == 1 {
		return rootNodes, nil
	} else {
		return nil, errors.New("Dummy")
	}

}

func processSoraNode(sora RootNode) error {
	for _, node := range sora.Fields {
		if node.Name == "connections" {
			processConnectionsNode(node)
		} else if node.Name == "stats" {
			processStatsNode(node)
		} else {
			// nop
		}
	}
	return nil
}

func processConnectionsNode(connections Node) error {
	debugPrint(connections)
	debugPrint(connections.Fields)

	for _, f := range connections.Fields {
		debugPrint("-------------------")
		debugPrint(f)
		debugPrint(f.Name)
		debugPrint(f.Fields)
		for _, f := range f.Fields {
			debugPrint("===================")
			debugPrint(f)
			debugPrint(f.Name)
			debugPrint(f.Fields)
		}
	}

	return nil
}

func processStatsNode(Node) error {
	// TODO: NYI
	return nil
}

// {
//   "_id": "a3c0a2d0-c4f2-11e7-b277-79c0643bd2c8",
//   "_type": "visualization",
//   "_source": {
//     "title": "Sora ongoing connections (TODO: host またぎ)",
//     "visState": "{
//        \"title\":\"Sora ongoing connections (TODO: host またぎ)\",
//        \"type\":\"metrics\",
//        \"params\":{
//           \"id\":\"61ca57f0-469d-11e7-af02-69e470af7417\",
//           \"type\":\"timeseries\",
//           \"series\":[{
//                  \"id\":\"61ca57f1-469d-11e7-af02-69e470af7417\",
//                  \"color\":\"#68BC00\",
//                  \"split_mode\":\"everything\",
//                  \"metrics\":[{
//                      \"id\":\"61ca57f2-469d-11e7-af02-69e470af7417\",
//                      \"type\":\"max\",
//                      \"field\":\"sora.stats.total_ongoing_connections\"}],
//                  \"seperate_axis\":0,
//                  \"axis_position\":\"right\",
//                  \"formatter\":\"number\",
//                  \"chart_type\":\"line\",
//                  \"line_width\":1,
//                  \"point_size\":1,
//                  \"fill\":0.5,
//                  \"stacked\":\"none\",
//                  \"label\":\"ongoing_connections\"}], // end of series
//        \"time_field\":\"@timestamp\",
//        \"index_pattern\":\"*\",
//        \"interval\":\"auto\",
//        \"axis_position\":\"left\",
//        \"axis_formatter\":\"number\",
//        \"show_legend\":1,
//        \"show_grid\":1,
//        \"axis_min\":\"0\"},
//        \"aggs\":[]}",  // end of visState
//     // _source のフィールド
//     "uiStateJSON": "{}",
//     "description": "",
//     "version": 1,
//     "kibanaSavedObjectMeta": {
//       "searchSourceJSON": "{}"
//     } // end of kibanaSavedObjectMeta
//   }  // end of _source
// },

// {
//     "title": "Sora total bytes",
//     "type": "metrics",
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
//             },
//             {
//                 "axis_position": "right",
//                 "chart_type": "line",
//                 "color": "#68BC00",
//                 "fill": "0",
//                 "formatter": "bytes",
//                 "id": "61ca57f1-469d-11e7-af02-69e470af7417",
//                 "label": "recieved",
//                 "line_width": "2",
//                 "metrics": [
//                     {
//                         "field": "sora.connections.rtp.total_received_bytes",
//                         "id": "61ca57f2-469d-11e7-af02-69e470af7417",
//                         "type": "max"
//                     },
//                     {
//                         "field": "61ca57f2-469d-11e7-af02-69e470af7417",
//                         "id": "b60bfb30-bfaf-11e7-ba99-7dd83649120a",
//                         "type": "derivative",
//                         "unit": "1s"
//                     }
//                 ],
//                 "point_size": "2",
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
//                 "terms_order_by": "61ca57f2-469d-11e7-af02-69e470af7417",
//                 "terms_size": "10",
//                 "value_template": "{{value}}/s"
//             }
//         ],
//         "show_grid": 1,
//         "show_legend": 1,
//         "time_field": "@timestamp",
//         "type": "timeseries"
//     }
//     "aggs": []
// }

func visualizationJson(
	prefix string, item string,
	splitMode string, derivative bool, termsField string,
	metricsType string,
	formatter string, axis_min int32) ([]byte, error) {
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
	return json.Marshal(values)
}

func jsonObj() map[string]interface{} {
	return make(map[string]interface{})
}

func print(arg []byte) {
	fmt.Printf("%s\n", arg)
}

func debugPrintf(format string, args ...interface{}) {
	fmt.Printf(format + "\n", args)
}

func debugPrint(arg interface{}) {
	fmt.Printf("%#v\n", arg)
}
