package optimizer

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	checker2 "github.com/ryan-berger/jpl/internal/ast/types/checker"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test-cp-*
var constantProp embed.FS

func readConstantPropTests() (map[string]string, map[string]string, error) {
	testMap := make(map[string]string)
	outputMap := make(map[string]string)

	e := fs.WalkDir(constantProp, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		f, err := constantProp.ReadFile(path)
		key := path[:len("testdata/test-cp-x")]

		if err != nil {
			return err
		}

		if strings.Contains(path, "output") {
			outputMap[key] = string(f)
		} else {
			testMap[key] = string(f)
		}

		return nil
	})
	if e != nil {
		return nil, nil, e
	}

	return testMap, outputMap, nil
}

func TestConstantProp(t *testing.T) {
	tests, outputs, _ := readConstantPropTests()
	for k, v := range tests {
		t.Run(k, func(t2 *testing.T) {
			lex, ok := lexer.Lex(v)
			assert.True(t2, ok)
			p, err := parser.Parse(lex)
			assert.NoError(t2, err, v)

			p = ConstantProp(p)

			p, _, err = checker2.Check(p)

			fmt.Println(p.String())

			assert.NoError(t2, err, v)
			assert.Equal(t2, outputs[k], p.SExpr())
		})

	}
}
