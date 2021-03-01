package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolTable_Copy(t *testing.T) {
	table1 := NewSymbolTable()
	table1["test"] = &Identifier{Type: Boolean}
	table1["x"] = &Identifier{Type: Boolean}
	table2 := table1.Copy()
	table2["y"] = &Identifier{Type: Boolean}

	fmt.Println(table1, table2)

	assert.NotEqual(t, table1, table2)
}
