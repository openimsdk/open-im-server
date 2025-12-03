package main

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"

	"github.com/likexian/gokit/assert"
	"gopkg.in/yaml.v3"
)

func Test_main(t *testing.T) {
	sourceYaml := ` # See the OWNERS docs at https://go.k8s.io/owners
approvers:
- dep-approvers
- thockin         # Network
- liggitt

labels:
- sig/architecture
`

	outputYaml := `# See the OWNERS docs at https://go.k8s.io/owners
approvers:
  - dep-approvers
  - thockin # Network
  - liggitt
labels:
  - sig/architecture
`
	node, _ := fetchYaml([]byte(sourceYaml))
	var output bytes.Buffer
	indent := 2
	writer := bufio.NewWriter(&output)
	_ = streamYaml(writer, &indent, node)
	_ = writer.Flush()
	assert.Equal(t, outputYaml, string(output.Bytes()), "yaml was not formatted correctly")
}

func Test_fetchYaml(t *testing.T) {
	type args struct {
		sourceYaml []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *yaml.Node
		wantErr bool
	}{
		{
			name: "Valid YAML",
			args: args{sourceYaml: []byte("key: value")},
			want: &yaml.Node{
				Kind:  yaml.MappingNode,
				Tag:   "!!map",
				Value: "",
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "key",
					},
					{
						Kind:  yaml.ScalarNode,
						Tag:   "!!str",
						Value: "value",
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid YAML",
			args:    args{sourceYaml: []byte("key:")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetchYaml(tt.args.sourceYaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fetchYaml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_streamYaml(t *testing.T) {
	type args struct {
		indent *int
		in     *yaml.Node
	}
	defaultIndent := 2
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "Valid YAML node with default indent",
			args: args{
				indent: &defaultIndent,
				in: &yaml.Node{
					Kind:  yaml.MappingNode,
					Tag:   "!!map",
					Value: "",
					Content: []*yaml.Node{
						{
							Kind:  yaml.ScalarNode,
							Tag:   "!!str",
							Value: "key",
						},
						{
							Kind:  yaml.ScalarNode,
							Tag:   "!!str",
							Value: "value",
						},
					},
				},
			},
			wantWriter: "key: value\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := streamYaml(writer, tt.args.indent, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("streamYaml() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("streamYaml() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
