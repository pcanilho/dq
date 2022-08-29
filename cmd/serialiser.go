package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/r3labs/diff/v3"
	"gopkg.in/yaml.v3"
	"strings"
)

type Serialiser interface {
	Serialise(any) string
}

type jsonSerialiser struct{}

func (j *jsonSerialiser) Serialise(i any) string {
	c, _ := json.MarshalIndent(i, "", " ")
	return string(c)
}

func newJSONSerialiser() *jsonSerialiser {
	return new(jsonSerialiser)
}

type yamlSerialiser struct{}

func newYAMLSerialiser() *yamlSerialiser {
	return new(yamlSerialiser)
}

func (y *yamlSerialiser) Serialise(i any) string {
	c, _ := yaml.Marshal(i)
	return fmt.Sprintf("---\n%s", string(c))
}

type junosSerialiser struct{}

func newJunosSerialiser() *junosSerialiser {
	return new(junosSerialiser)
}

func (j *junosSerialiser) Serialise(it any) string {
	settings, ok := it.(map[string]diff.Changelog)
	if !ok {
		return ""
	}

	var sb strings.Builder
	for header, changelog := range settings {
		sb.WriteString(fmt.Sprintf("\n[%s]\n", header))
		for _, change := range changelog {
			removal := color.New(color.FgRed)
			addition := color.New(color.FgGreen)
			var opType *color.Color
			switch change.Type {
			case diff.UPDATE:
				opType = color.New(color.FgHiYellow, color.Bold)
			case diff.DELETE:

				opType = color.New(color.FgHiRed, color.Bold)
			}

			sb.WriteString(fmt.Sprintf("[%s] (%s)\n", strings.Join(change.Path, "."), opType.Sprint(change.Type)))
			sb.WriteString(removal.Sprintf("-\t%v\n", change.From))
			sb.WriteString(addition.Sprintf("+\t%v\n", change.To))
		}
	}
	return sb.String()
}
