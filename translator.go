package translator

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type Translator interface {
	Translate(langCode, key string, args ...map[string]any) string
}

type translator struct {
	dir        string
	fileType   string
	translates map[string]map[string]string
}

const (
	Json = "json"
	Yaml = "yaml"
	Toml = "toml"
)

func New(dir, fileType string) Translator {
	t := &translator{
		dir:        dir,
		fileType:   fileType,
		translates: make(map[string]map[string]string),
	}
	if len(t.dir) == 0 {
		return t
	}
	if _, err := os.Stat(t.dir); os.IsNotExist(err) {
		panic(ErrorInvalidDir)
	}
	err := t.walk()
	if err != nil {
		panic(err)
	}
	return t
}

func (t *translator) Translate(langCode, key string, args ...map[string]any) string {
	langTranslates, ok := t.translates[langCode]
	if !ok {
		return key
	}
	translate, ok := langTranslates[key]
	if !ok {
		return key
	}
	translate = replaceArgs(translate, args...)
	return translate
}

func (t *translator) walk() error {
	if err := filepath.Walk(
		t.dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(info.Name(), "."+t.fileType) {
				return nil
			}
			lang := strings.TrimSuffix(info.Name(), "."+t.fileType)
			if t.translates[lang] == nil {
				t.translates[lang] = make(map[string]string)
			}
			dir := strings.TrimPrefix(t.dir, "./")
			subpath := strings.TrimPrefix(strings.TrimSuffix(path, info.Name()), dir)
			subpath = strings.TrimPrefix(subpath, "/")
			subpath = strings.TrimSuffix(subpath, "/")
			if err := t.read(lang, path, createKeyPrefixFromPath(subpath)); err != nil {
				return err
			}
			return nil
		},
	); err != nil {
		return err
	}
	return nil
}

func (t *translator) read(lang, path, prefix string) error {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fileData := make(map[string]any)
	switch t.fileType {
	case Json:
		if err := json.Unmarshal(fileBytes, &fileData); err != nil {
			return err
		}
	case Toml:
		if err := toml.Unmarshal(fileBytes, &fileData); err != nil {
			return err
		}
	case Yaml:
		if err := yaml.Unmarshal(fileBytes, &fileData); err != nil {
			return err
		}
	}
	return t.parse(lang, prefix, fileData)
}

func (t *translator) parse(lang, prefix string, data map[string]any) error {
	hasPrefix := len(prefix) > 0
	for dataKey, item := range data {
		var key string
		if hasPrefix {
			key = fmt.Sprintf("%s.%v", prefix, dataKey)
		}
		if !hasPrefix {
			key = fmt.Sprintf("%v", dataKey)
		}
		subdata, ok := item.(map[string]any)
		if !ok {
			t.translates[lang][key] = fmt.Sprintf("%v", item)
		}
		if ok {
			return t.parse(lang, key, subdata)
		}
	}
	return nil
}
