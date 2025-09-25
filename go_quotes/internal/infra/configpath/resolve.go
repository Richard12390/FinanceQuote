package configpath

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Resolve(cliPath string) string {

	// Config path
	if path := strings.TrimSpace(os.Getenv(EnvConfigPath)); path != "" {
		if absPath := checkPath(path); absPath != "" {
			return absPath
		}
		return ""
	}

	// Environment variable
	if env := strings.TrimSpace(os.Getenv(EnvConfigPath)); env != "" {
		if abs := checkPath(env); abs != "" {
			return abs
		}
	}

	// Default folders
	for _, dir := range []string{
		filepath.Join("configs", "pipelines"),
		filepath.Join("configs", "pipeline"),
	} {
		fmt.Println(filepath.Join("configs", "pipelines"))
		if absPath := checkPath(dir); absPath != "" {
			return absPath
		}
	}

	// Default files
	for _, file := range []string{
		"config.yml",
		"config.yaml",
		filepath.Join("configs", "config.yml"),
		filepath.Join("configs", "config.yaml"),
	} {
		if absPath := checkPath(file); absPath != "" {
			return absPath
		}
	}
	return ""
}
func checkPath(path string) string {
	path = expand(path)
	st, err := os.Stat(path)
	if err != nil {
		return ""
	}
	if st.IsDir() {
		if !hasYaml(path) {
			return ""
		}
	}
	abs, _ := filepath.Abs(path)
	return abs
}

func hasYaml(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
			return true
		}
	}
	return false
}

func expand(path string) string {
	path = os.ExpandEnv(path)
	if strings.HasPrefix(path, "~") {
		if home, err := os.UserHomeDir(); err == nil {
			if path == "~" {
				path = home
			} else if len(path) >= 2 && (path[1] == '/' || path[1] == '\\') {
				path = filepath.Join(home, path[2:])
			}
		}
	}
	return filepath.Clean(path)
}
