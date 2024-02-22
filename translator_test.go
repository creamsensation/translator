package translator

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTranslator(t *testing.T) {
	rootDir := t.TempDir()
	translatesDir := fmt.Sprintf("%s/translates", rootDir)
	assert.Nil(t, os.MkdirAll(translatesDir, os.ModePerm))
	t.Run(
		"simple", func(t *testing.T) {
			enPath := fmt.Sprintf("%s/en.yaml", translatesDir)
			assert.Nil(
				t,
				os.WriteFile(
					enPath,
					[]byte("label:\n  email: E-mail\n  name: Name"),
					os.ModePerm,
				),
			)
			translator := New(Config{translatesDir, Yaml})
			assert.Equal(t, "E-mail", translator.Translate("en", "label.email"))
			assert.Equal(t, "Name", translator.Translate("en", "label.name"))
		},
	)

	t.Run(
		"nested", func(t *testing.T) {
			assert.Nil(t, os.MkdirAll(fmt.Sprintf("%s/global/label", translatesDir), os.ModePerm))
			enPath := fmt.Sprintf("%s/global/label/en.yaml", translatesDir)
			assert.Nil(
				t,
				os.WriteFile(
					enPath,
					[]byte("email: E-mail\nname: Name"),
					os.ModePerm,
				),
			)
			translator := New(Config{translatesDir, Yaml})
			assert.Equal(t, "E-mail", translator.Translate("en", "global.label.email"))
			assert.Equal(t, "Name", translator.Translate("en", "global.label.name"))
		},
	)
}
