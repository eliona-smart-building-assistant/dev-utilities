package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Provide path to directory with asset type defintions.")
		return
	}
	dir := os.Args[1]
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), "asset-type-") && strings.HasSuffix(file.Name(), ".json") {
			assetType := strings.TrimSuffix(strings.TrimPrefix(file.Name(), "asset-type-"), ".json")
			processFile(dir, file.Name(), unexport(kebabToCamelCase(assetType)))
		}
	}
}

type Attribute struct {
	Enable  bool   `json:"enable"`
	Name    string `json:"name"`
	Subtype string `json:"subtype"`
	Type    string `json:"type"`
}

type Input struct {
	Attributes []Attribute `json:"attributes"`
}

func processFile(dir, filename, assetType string) {
	data, err := ioutil.ReadFile(filepath.Join(dir, filename))
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", filename, err)
		return
	}

	var input Input
	err = json.Unmarshal(data, &input)
	if err != nil {
		fmt.Printf("Error unmarshalling JSON from file %s: %v\n", filename, err)
		return
	}

	infoDataPayload := map[string]string{}
	statusDataPayload := map[string]string{}
	inputDataPayload := map[string]string{}

	for _, attr := range input.Attributes {
		if !attr.Enable {
			continue
		}

		fieldName := snakeToCamelCase(attr.Name)

		switch attr.Subtype {
		case "info":
			infoDataPayload[fieldName] = attr.Name
		case "status":
			statusDataPayload[fieldName] = attr.Name
		case "input":
			inputDataPayload[fieldName] = attr.Name
		default:
			fmt.Printf("%s: unknown subtype %s", filename, attr.Subtype)
		}
	}

	printStruct(assetType+"InfoDataPayload", infoDataPayload)
	printStruct(assetType+"StatusDataPayload", statusDataPayload)
	printStruct(assetType+"InputDataPayload", inputDataPayload)
}

func printStruct(name string, fields map[string]string) {
	if len(fields) == 0 {
		return
	}
	fmt.Printf("type %s struct {\n", name)
	for fieldName, jsonFieldName := range fields {
		fmt.Printf("\t%s string `json:\"%s\"`\n", fieldName, jsonFieldName)
	}
	fmt.Print("}\n\n")
}

func snakeToCamelCase(s string) string {
	words := strings.Split(s, "_")
	result := ""
	for _, word := range words {
		result += strings.Title(word)
	}
	return result
}

func kebabToCamelCase(s string) string {
	words := strings.Split(s, "-")
	result := ""
	for _, word := range words {
		result += strings.Title(word)
	}
	return result
}

func unexport(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := strings.ToLower(string(s[0]))
	restOfString := s[1:]
	return firstChar + restOfString
}
