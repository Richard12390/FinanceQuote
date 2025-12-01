package conf

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Load(path string) (Config, error) {
	var out Config
	st, err := os.Stat(path)
	if err != nil {
		return out, err
	}

	if st.IsDir() {
		var pipelines []Pipeline
		walkErr := filepath.WalkDir(path, func(p string, d fs.DirEntry, walkErr error) error {
			if walkErr != nil || d.IsDir() || !isYAML(p) {
				return walkErr
			}
			raw, err := os.ReadFile(p)
			if err != nil {
				return err
			}

			var cfg Config
			if yaml.Unmarshal(raw, &cfg) == nil {
				if cfg.NATSUrl != "" && out.NATSUrl == "" {
					out.NATSUrl = cfg.NATSUrl
				}
				if len(cfg.Pipelines) > 0 {
					pipelines = append(pipelines, cfg.Pipelines...)
					return nil
				}
				if cfg.NATSUrl != "" && out.NATSUrl == cfg.NATSUrl {
					return nil
				}
			}

			var pl Pipeline
			if err := yaml.Unmarshal(raw, &pl); err != nil {
				return err
			}
			if pl.Name == "" {
				return fmt.Errorf("pipeline file missing name: %s", p)
			}
			pipelines = append(pipelines, pl)
			return nil
		})
		if walkErr != nil {
			return out, walkErr
		}
		out.Pipelines = pipelines
		return out, nil
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := yaml.Unmarshal(raw, &out); err != nil {
		return out, err
	}
	return out, nil
}
func isYAML(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yml" || ext == ".yaml"
}
