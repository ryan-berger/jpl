package optimizer

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/ryan-berger/jpl/internal/ast/types/checker"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test-dce-*
var deadCodeTests embed.FS

func readDeadCodeTests() (map[string]string, map[string]string, error) {
	testMap := make(map[string]string)
	outputMap := make(map[string]string)

	e := fs.WalkDir(deadCodeTests, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil {
			return err
		}

		f, err := deadCodeTests.ReadFile(path)
		key := path[:len("testdata/test-dce-x")]

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

func TestDeadCode(t *testing.T) {
	tests, output, _ := readDeadCodeTests()
	for k, v := range tests {
		t.Run(k, func(tes *testing.T) {
			lex, ok := lexer.Lex(v)
			assert.True(tes, ok)
			p, err := parser.Parse(lex)
			assert.NoError(tes, err)

			fmt.Println(p.SExpr())

			p = DeadCode(p)
			p, _, err = checker.Check(p)
			assert.NoError(tes, err, v)
			assert.Equal(tes, output[k], p.SExpr())
		})
	}
}
