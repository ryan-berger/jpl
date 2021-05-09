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

//go:embed testdata/test-peep-*
var peepTests embed.FS

func readPeepTests() (map[string]string, map[string]string, error){
	testMap := make(map[string]string)
	outputMap := make(map[string]string)


	e := fs.WalkDir(peepTests, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if err != nil  {
			return err
		}

		f, err := peepTests.ReadFile(path)
		key := path[:len("testdata/test-peep-x")]

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

func TestPeephole(t *testing.T) {
	tests, _, _ := readPeepTests()
	for _, v := range tests {
		lex, ok := lexer.Lex(v)
		assert.True(t, ok)
		p, err := parser.Parse(lex)
		assert.NoError(t, err)
		p = Peephole(p)
		p, _, err = checker2.Check(p)
		assert.NoError(t, err, v)
		fmt.Println(p.String())
	}
}
