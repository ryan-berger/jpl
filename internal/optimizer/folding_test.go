package optimizer

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/ryan-berger/jpl/internal/ast/typed"
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/parser"
	"github.com/stretchr/testify/assert"
)

//go:embed testdata/test-cf-*
var foldTests embed.FS

func readTests() (map[string]string, map[string]string, error){
	testMap := make(map[string]string)
	outputMap := make(map[string]string)


	e := fs.WalkDir(foldTests, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil  {
			return err
		}

		f, err := foldTests.ReadFile(path)
		key := path[:len("testdata/test-cf-x")]

		if err != nil {
			return err
		}

		if strings.Contains(path,"output") {
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
	tests, _, _ := readTests()
	for _, v := range tests {
		lex, ok := lexer.Lex(v)
		assert.True(t, ok)
		p, err := parser.Parse(lex)
		assert.NoError(t, err)
		p = ConstantFold(p)
		p, _, err = typed.Check(p)
		assert.NoError(t, err)
		fmt.Println(p.SExpr())
	}
}
