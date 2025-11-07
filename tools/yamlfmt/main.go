// OPENIM plan on prow tools
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	// Prow OWNERs file defines the default indent as 2 spaces.
	indent := flag.Int("indent", 2, "default indent")
	flag.Parse()
	for _, path := range flag.Args() {
		sourceYaml, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			continue
		}
		rootNode, err := fetchYaml(sourceYaml)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			continue
		}
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			continue
		}
		err = streamYaml(writer, indent, rootNode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
			continue
		}
	}
}

func fetchYaml(sourceYaml []byte) (*yaml.Node, error) {
	rootNode := yaml.Node{}
	err := yaml.Unmarshal(sourceYaml, &rootNode)
	if err != nil {
		return nil, err
	}
	return &rootNode, nil
}

func streamYaml(writer io.Writer, indent *int, in *yaml.Node) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(*indent)
	err := encoder.Encode(in)
	if err != nil {
		return err
	}
	return encoder.Close()
}
