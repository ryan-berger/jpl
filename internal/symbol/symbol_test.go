package symbol

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/types"
)

func TestSymbolTable_Copy(t *testing.T) {
	table1 := NewSymbolTable()
	table1["test"] = &Identifier{Type: types.Boolean}
	table1["x"] = &Identifier{Type: types.Boolean}
	table2 := table1.Copy()
	table2["y"] = &Identifier{Type: types.Boolean}

	fmt.Println(table1, table2)

	assert.NotEqual(t, table1, table2)
}
