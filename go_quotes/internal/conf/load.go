package conf

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	var out Config
	st, err := os.Stat(path)
	if err != nil {
		return out, err
	}

	if st.IsDir() {
		var pls []Pipeline
		walkErr := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			ext := filepath.Ext(p)
			if ext != ".yml" && ext != ".yaml" {
				return nil
			}
			b, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			var singlepl struct {
				Pipelines []Pipeline `yaml:"pipelines"`
			}
			if yaml.Unmarshal(b, &singlepl) == nil && len(singlepl.Pipelines) > 0 {
				pls = append(pls, singlepl.Pipelines...)
				return nil
			}
			var multipls Pipeline
			if err := yaml.Unmarshal(b, &multipls); err != nil {
				return err
			}
			if multipls.Name == "" {
				return errors.New("pipeline file missing name: " + p)
			}
			pls = append(pls, multipls)
			return nil
		})
		if walkErr != nil {
			return out, walkErr
		}
		out.Pipelines = pls
		return out, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return out, err
	}
	return out, nil
}
