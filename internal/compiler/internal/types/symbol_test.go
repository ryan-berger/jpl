package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTable_Copy(t *testing.T) {
	table1 := NewSymbolTable()
	table1["test"] = &Identifier{Type: boolean}
	table1["x"] = &Identifier{Type: boolean}
	table2 := table1.Copy()
	table2["y"] = &Identifier{Type: boolean}

	fmt.Println(table1, table2)

	assert.NotEqual(t, table1, table2)
}
