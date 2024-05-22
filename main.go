package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

const tagsRegex = `\[\[(.*?)\]\]`

type Tag struct {
	Name string
}

type Source struct {
	Name string
	Path string
	Tags []Tag
}

type Node struct {
	Name    string
	Links   *[]Node
	Sources *[]Source
}

func readFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func CleanTag(tag string) string {
	tag = tag[2 : len(tag)-2]
	tag = strings.TrimSpace(tag)
	return tag
}

func ExtractMapKeys(m map[string]bool) []string {
	keys := reflect.ValueOf(m).MapKeys()
	convKeys := make([]string, len(keys))
	for _, key := range keys {
		convKeys = append(convKeys, key.String())
	}
	return convKeys
}

func GetTagsFromString(tagsRegex, s string) []string {
	return regexp.MustCompile(tagsRegex).FindAllString(s, -1)
}

func ExtractTagsFromString(fileContent []byte) []string {
	tags := make(map[string]bool, 0)
	matches := GetTagsFromString(tagsRegex, string(fileContent))

	for _, match := range matches {
		tag := CleanTag(match)
		if !tags[tag] {
			tags[tag] = true
		}
	}

	return ExtractMapKeys(tags)
}

func ProccessSource(source string) *Source {
	exec, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd, err := filepath.EvalSymlinks(exec)
	if err != nil {
		panic(err)
	}

	filePath := filepath.Join(wd, source)
	fileContent, error := readFile(filePath)
	if error != nil {
		panic(error)
	}

	tags := []Tag{}
	for _, tag := range ExtractTagsFromString(fileContent) {
		tags = append(tags, Tag{Name: tag})
	}

	return &Source{
		Name: source,
		Path: filePath,
		Tags: tags,
	}
}

func main() {
	file := os.Args[1]

	source := ProccessSource(file)

	getTagHash := func(t Tag) string {
		return t.Name
	}

	g := graph.New(getTagHash)

	for _, tag := range source.Tags {
		g.AddVertex(tag)
	}

	fmt.Printf("Source: %+v\n", source)

	output, _ := os.Create("./mygraph.gv")
	draw.DOT(g, output)
}
