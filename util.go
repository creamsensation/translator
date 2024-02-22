package translator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func createKeyPrefixFromPath(path string) string {
	if len(path) == 0 {
		return ""
	}
	result := make([]string, 0)
	for _, p := range strings.Split(path, "/") {
		if len(p) == 0 {
			continue
		}
		result = append(result, strcase.ToKebab(p))
	}
	return strings.Join(result, ".")
}

func replaceArgs(value string, args ...map[string]any) string {
	if len(args) == 0 || !strings.Contains(value, "{{") || !strings.Contains(value, "}}") {
		return value
	}
	value = strings.Replace(value, "{{ ", "{{", -1)
	value = strings.Replace(value, " }}", "}}", -1)
	for ak, av := range args[0] {
		value = strings.Replace(value, fmt.Sprintf("{{%s}}", ak), fmt.Sprintf("%v", av), -1)
	}
	return value
}
