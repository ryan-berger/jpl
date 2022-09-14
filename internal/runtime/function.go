package runtime

import "github.com/ryan-berger/jpl/internal/ast/types"

type Param struct {
	Name string
	Type types.Type
}

type Function struct {
	Name   string
	Return types.Type
	Params []Param
}

var fns = []Function{
	{
		Name:   "get_time",
		Return: types.Float,
	},
	{
		Name: "print",
		Params: []Param{
			{
				Name: "val",
				Type: String,
			},
		},
		Return: Void,
	},
	{
		Name: "show",
		Params: []Param{
			{
				Name: "type",
				Type: String,
			},
			{
				Name: "val",
				Type: &Pointer{
					Inner: Opaque,
				},
			},
		},
		Return: Void,
	},
	{
		Name: "read_image",
		Params: []Param{
			{
				Name: "res",
				Type: &Pointer{
					Inner: types.Pict,
				},
			},
			{
				Name: "file_name",
				Type: String,
			},
		},
		Return: Void,
	},
	{
		Name: "write_image",
		Params: []Param{
			{
				Name: "pict",
				Type: &Pointer{
					Inner:     types.Pict,
					Alignment: 8,
				},
			},
			{
				Name: "file_name",
				Type: String,
			},
		},
		Return: Void,
	},
	{
		Name: "fail_assertion",
		Params: []Param{
			{
				Name: "message",
				Type: String,
			},
		},
		Return: Void,
	},
}

var Functions = make(map[string]Function)

func init() {
	for _, fn := range fns {
		Functions[fn.Name] = fn
	}
}
