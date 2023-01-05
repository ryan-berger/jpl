package optimizer

import (
	"embed"
	"io/fs"
	"strings"
	"testing"

	checker2 "github.com/ryan-berger/jpl/internal/ast/types/checker"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test-cf-*
var foldTests embed.FS

func readTests() (map[string]string, map[string]string, error) {
	testMap := make(map[string]string)
	outputMap := make(map[string]string)

	e := fs.WalkDir(foldTests, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		f, err := foldTests.ReadFile(path)
		key := path[:len("testdata/test-cf-x")]

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

func TestConstantFolding(t *testing.T) {
	tests, outputs, _ := readTests()
	for k, v := range tests {
		t.Run(k, func(t2 *testing.T) {
			lex, ok := lexer.Lex(v)

			assert.True(t2, ok)

			p, err := parser.Parse(lex)
			assert.NoError(t2, err)
			p = ConstantFold(p)

			p, _, err = checker2.Check(p)

			assert.NoError(t2, err, v)
			assert.Equal(t2, outputs[k], p.SExpr())
		})

	}
}
