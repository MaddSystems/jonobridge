package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name             string     `yaml:"name" json:"name"`
	Version          string     `yaml:"version" json:"version"`
	Description      string     `yaml:"description" json:"description"`
	GruleContextName string     `yaml:"grule_context_name" json:"grule_context_name"`
	Functions        []Function `yaml:"functions" json:"functions"`
}

type Function struct {
	Name        string      `yaml:"name" json:"name"`
	Description string      `yaml:"description" json:"description"`
	Parameters  []Parameter `yaml:"parameters" json:"parameters"`
	ReturnType  string      `yaml:"return_type,omitempty" json:"return_type,omitempty"`
	Example     string      `yaml:"example,omitempty" json:"example,omitempty"`
}

type Parameter struct {
	Name string `yaml:"name" json:"name"`
	Type string `yaml:"type" json:"type"`
}

func GenerateFromManifests(capabilitiesDir string) ([]byte, error) {
	manifests := []Manifest{}

	err := filepath.Walk(capabilitiesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "manifest.yaml" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			var m Manifest
			if err := yaml.Unmarshal(data, &m); err != nil {
				return fmt.Errorf("error parsing %s: %v", path, err)
			}
			manifests = append(manifests, m)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	output := map[string]interface{}{
		"version":      "1.0.0",
		"capabilities": manifests,
	}

	return json.MarshalIndent(output, "", "  ")
}
