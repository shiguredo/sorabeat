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
	splitMode := "everything"
	metricsType := "max"
	formatter := "byte"
	var axis_min int32
	axis_min = 0
	vres, reserr := visualizationJson(prefix, item, splitMode, metricsType, formatter, axis_min)
	if reserr != nil {
		debugPrint(reserr)
		os.Exit(3)
	}
	print(vres)

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

func visualizationJson(prefix string, item string, splitMode string, metricsType string,
	formatter string, axis_min int32) ([]byte, error) {
	values := jsonObj()
	{
		values["_id"] = uuid.NewV4()
		values["_type"] = "visualization"
		source := jsonObj()
		{
			source["title"] = "[Sora] " + item
			visState := jsonObj()
			{
				visState["title"] = "[Sora] " + item
				visState["type"] = "metrics"
				idParams := uuid.NewV4()
				params := jsonObj()
				{
					params["id"] = idParams
					params["type"] = "timeseries"
					var series [1]interface{}
					series0 := jsonObj()
					{
						series0["id"] = idParams
						series0["color"] = "#68BC00"
						// TODO: term で split する場合
						series0["split_mode"] = splitMode
						metrics := jsonObj()
						{
							metrics["id"] = idParams
							metrics["type"] = metricsType // e.g. "max"
							metrics["field"] = prefix + item
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
					series[0] = series0
					params["series"] = series
				}
				visState["params"] = params
				visState["time_field"] = "@timestamp"
				visState["index_pattern"] = "*"
				visState["interval"] = "auto"
				visState["axis_position"] = "left"
				visState["axis_formatter"] = formatter
				visState["show_legend"] = 1
				visState["show_grid"] = 1
				visState["axis_min"] = axis_min
				visState["aggs"] = []string{}
			}
			visStateBytes, _ := json.Marshal(visState)
			source["visState"] = string(visStateBytes[:])
		}
		values["_source"] = source
		values["uiStateJson"] = "{}"
		values["description"] = "[Sora]" + item
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
