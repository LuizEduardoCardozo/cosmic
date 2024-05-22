package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
)

const tagsRegex = `\[\[.*?\]\]`

type Tag struct {
	Name string
}

type Source struct {
	Name string
	Path string
	Tags []Tag
}

var Sources = make(map[string]Source, 0)

func HashSource(s Source) (bs string) {
	h := sha256.New()
	h.Write([]byte(s.Path))
	bs = string(h.Sum(nil))
	return
}

func readFile(filename string) ([]byte, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func CleanTag(tag string) string {
	return strings.TrimSpace(tag[2 : len(tag)-2])
}

func ExtractMapKeys(m map[string]bool) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetTagsFromString(tagsRegex, s string) []string {
	tags := regexp.MustCompile(tagsRegex).FindAllString(s, -1)
	return tags
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

func ProccessSource(sources []string) {
	exec, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	wd, err := filepath.EvalSymlinks(exec)
	if err != nil {
		panic(err)
	}

	for _, source := range sources {
		filePath := filepath.Join(wd, source)
		fileContent, error := readFile(filePath)
		if error != nil {
			panic(error)
		}

		tags := make([]Tag, 0)
		for _, tag := range ExtractTagsFromString(fileContent) {
			tags = append(tags, Tag{Name: tag})
		}

		s := Source{
			Name: source,
			Path: filePath,
			Tags: tags,
		}

		hash := HashSource(s)

		Sources[hash] = s
	}
}

func DrawGraph() {
	g := graph.New(func(s string) string { return s })

	for _, source := range Sources {
		g.AddVertex(source.Name,
			graph.VertexAttribute("colorscheme", "blues3"),
			graph.VertexAttribute("style", "filled"),
			graph.VertexAttribute("color", "2"),
			graph.VertexAttribute("fillcolor", "1"),
		)
		for _, tag := range source.Tags {
			g.AddVertex(tag.Name,
				graph.VertexAttribute("colorscheme", "greens3"),
				graph.VertexAttribute("style", "filled"),
				graph.VertexAttribute("color", "2"),
				graph.VertexAttribute("fillcolor", "1"),
			)
			g.AddEdge(source.Name, tag.Name)
		}
		fmt.Printf("Source: %+v\n", source)
	}

	output, _ := os.Create("./mygraph.gv")
	draw.DOT(g, output)
}

func main() {
	file := os.Args[1:3]
	fmt.Println("Proccessing filess: ", file)

	ProccessSource(file)
	DrawGraph()
}
