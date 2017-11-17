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
		fmt.Println(err1)
		os.Exit(1)
	}

	res, err2 := processSoraFields(buf)
	if err2 != nil {
		fmt.Println(err2)
		os.Exit(2)
	}
	fmt.Println(res)
	fmt.Println("SUCCEEDED!! ＼（＾ ＾）／")
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
	name string
	Title string
	Type string
	Description string
	Fields []Node `yaml:"fields,omitempty"`
}

func processSoraFields(buf []byte) (interface{}, error) {
	var contents []RootNode
	err1 := yaml.Unmarshal(buf, &contents)
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(contents[0])
	fmt.Printf("%#v\n", contents[0])
	fmt.Println(contents[0].Key)

	// if first, ok := contents[0].(map[string]interface{}); ok {
	// 	if key1, ok := first["key"]; ok {
	// 		fmt.Println("===========key1===========")
	// 		fmt.Println(key1)
	// 		fmt.Println("===========key1===========")
	// 	} else {
	// 		return nil, errors.New("Fail to fetch key 'key'")
	// 	}
	// } else {
	// 	return nil, errors.New("Fail to type conversion")
	// }

	if 1 == 1 {
		return contents, nil
	} else {
		return nil, errors.New("Dummy")
	}
}
