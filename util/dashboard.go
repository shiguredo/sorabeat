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
	"errors"
	"os"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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


func debugPrintf(format string, args ...interface{}) {
	fmt.Printf(format + "\n", args)
}

func debugPrint(arg interface{}) {
	fmt.Println(arg)
}
