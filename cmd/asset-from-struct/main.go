package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"reflect"
	"strings"
)

type attribute struct {
	Enable      bool              `json:"enable"`
	Name        string            `json:"name"`
	Subtype     string            `json:"subtype"`
	Type        string            `json:"type,omitempty"`
	Translation map[string]string `json:"translation,omitempty"`
	Unit        string            `json:"unit,omitempty"`
}

type assetTypeDef struct {
	Attributes  []attribute       `json:"attributes"`
	Custom      bool              `json:"custom"`
	Icon        string            `json:"icon"`
	Name        string            `json:"name"`
	Translation map[string]string `json:"translation"`
	Urldoc      string            `json:"urldoc"`
	Vendor      string            `json:"vendor"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter struct (finish with double newline or EOF):")

	var structStrBuilder strings.Builder
	previousLineWasEmpty := false
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// EOF might be a valid case to break out if it's the end of the input.
				break
			}
			panic(err)
		}

		// Check for a double newline, which indicates the end of the input
		if line == "\n" {
			if previousLineWasEmpty {
				break
			}
			previousLineWasEmpty = true
		} else {
			previousLineWasEmpty = false
		}

		structStrBuilder.WriteString(line)
	}
	structStr := structStrBuilder.String()

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", "package main\n"+structStr, 0)
	if err != nil {
		panic(err)
	}

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			assetType := processStruct(structType)
			jsonAssetType, _ := json.MarshalIndent(assetType, "", "\t")
			fmt.Println(string(jsonAssetType))
		}
	}
}

func processStruct(structType *ast.StructType) assetTypeDef {
	assetType := assetTypeDef{
		Custom: true,
		Icon:   "...",
		Name:   "...",
		Translation: map[string]string{
			"de": "...",
			"en": "...",
		},
		Urldoc: "...",
		Vendor: "...",
	}

	for _, field := range structType.Fields.List {
		if field.Names == nil || field.Tag == nil {
			continue
		}

		tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1]) // Remove the surrounding quotes.
		subtype := tag.Get("subtype")

		if subtype == "" {
			continue
		}

		elionaTag := tag.Get("eliona")
		elionaValues := strings.Split(elionaTag, ",")

		attr := attribute{
			Enable:  true,
			Name:    elionaValues[0], // Use the name part of the tag value (before the ",").
			Subtype: subtype,
			Translation: map[string]string{
				"de": "...",
				"en": "...",
			},
			Type: "...",
			Unit: "...",
		}

		assetType.Attributes = append(assetType.Attributes, attr)
	}

	return assetType
}
