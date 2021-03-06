package hclhil

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	hclast "github.com/hashicorp/hcl/hcl/ast"
	hcltoken "github.com/hashicorp/hcl/hcl/token"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-zcl/zcl"
)

func TestBodyPartialContent(t *testing.T) {
	tests := []struct {
		Source    string
		Schema    *zcl.BodySchema
		Want      *zcl.BodyContent
		DiagCount int
	}{
		{
			``,
			&zcl.BodySchema{},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				MissingItemRange: zcl.Range{
					Filename: "<unknown>",
				},
			},
			0,
		},
		{
			`foo = 1`,
			&zcl.BodySchema{},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			0,
		},
		{
			`foo = 1`,
			&zcl.BodySchema{
				Attributes: []zcl.AttributeSchema{
					{
						Name: "foo",
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{
					"foo": {
						Name: "foo",
						Expr: &expression{
							src: &hclast.LiteralType{
								Token: hcltoken.Token{
									Type: hcltoken.NUMBER,
									Text: `1`,
									Pos: hcltoken.Pos{
										Offset: 6,
										Line:   1,
										Column: 7,
									},
								},
							},
						},

						Range: zcl.Range{
							Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
							End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
						},
						NameRange: zcl.Range{
							Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
							End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
						},
					},
				},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			0,
		},
		{
			``,
			&zcl.BodySchema{
				Attributes: []zcl.AttributeSchema{
					{
						Name:     "foo",
						Required: true,
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				MissingItemRange: zcl.Range{
					Filename: "<unknown>",
				},
			},
			1, // missing required attribute
		},
		{
			`foo {}`,
			&zcl.BodySchema{
				Blocks: []zcl.BlockHeaderSchema{
					{
						Type: "foo",
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				Blocks: zcl.Blocks{
					{
						Type: "foo",
						Body: &body{
							oli: &hclast.ObjectList{},
						},
						DefRange: zcl.Range{
							Start: zcl.Pos{Byte: 4, Line: 1, Column: 5},
							End:   zcl.Pos{Byte: 5, Line: 1, Column: 6},
						},
						TypeRange: zcl.Range{
							Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
							End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
						},
					},
				},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			0,
		},
		{
			`foo "unwanted" {}`,
			&zcl.BodySchema{
				Blocks: []zcl.BlockHeaderSchema{
					{
						Type: "foo",
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			1, // no labels are expected
		},
		{
			`foo {}`,
			&zcl.BodySchema{
				Blocks: []zcl.BlockHeaderSchema{
					{
						Type:       "foo",
						LabelNames: []string{"name"},
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			1, // missing name
		},
		{
			`foo "wanted" {}`,
			&zcl.BodySchema{
				Blocks: []zcl.BlockHeaderSchema{
					{
						Type:       "foo",
						LabelNames: []string{"name"},
					},
				},
			},
			&zcl.BodyContent{
				Attributes: zcl.Attributes{},
				Blocks: zcl.Blocks{
					{
						Type:   "foo",
						Labels: []string{"wanted"},
						Body: &body{
							oli: &hclast.ObjectList{},
						},
						DefRange: zcl.Range{
							Start: zcl.Pos{Byte: 13, Line: 1, Column: 14},
							End:   zcl.Pos{Byte: 14, Line: 1, Column: 15},
						},
						TypeRange: zcl.Range{
							Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
							End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
						},
						LabelRanges: []zcl.Range{
							{
								Start: zcl.Pos{Byte: 4, Line: 1, Column: 5},
								End:   zcl.Pos{Byte: 5, Line: 1, Column: 6},
							},
						},
					},
				},
				MissingItemRange: zcl.Range{
					Start: zcl.Pos{Line: 1, Column: 1, Byte: 0},
					End:   zcl.Pos{Line: 1, Column: 2, Byte: 1},
				},
			},
			0,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			file, diags := Parse([]byte(test.Source), "test.hcl")
			if len(diags) != 0 {
				t.Fatalf("diagnostics from parse: %s", diags.Error())
			}

			got, _, diags := file.Body.PartialContent(test.Schema)
			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf(" - %s", diag.Error())
				}
			}

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", spew.Sdump(got), spew.Sdump(test.Want))
			}
		})
	}

}

func TestBodyJustAttributes(t *testing.T) {
	tests := []struct {
		Source    string
		Want      zcl.Attributes
		DiagCount int
	}{
		{
			``,
			zcl.Attributes{},
			0,
		},
		{
			`foo = "a"`,
			zcl.Attributes{
				"foo": &zcl.Attribute{
					Name: "foo",
					Expr: &expression{
						src: &hclast.LiteralType{
							Token: hcltoken.Token{
								Type: hcltoken.STRING,
								Pos: hcltoken.Pos{
									Offset: 6,
									Line:   1,
									Column: 7,
								},
								Text: `"a"`,
							},
						},
					},
					Range: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
					NameRange: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
				},
			},
			0,
		},
		{
			`foo = {}`,
			zcl.Attributes{
				"foo": &zcl.Attribute{
					Name: "foo",
					Expr: &expression{
						src: &hclast.ObjectType{
							List: &hclast.ObjectList{},
							Lbrace: hcltoken.Pos{
								Offset: 6,
								Line:   1,
								Column: 7,
							},
							Rbrace: hcltoken.Pos{
								Offset: 7,
								Line:   1,
								Column: 8,
							},
						},
					},
					Range: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
					NameRange: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
				},
			},
			0,
		},
		{
			`foo {}`,
			zcl.Attributes{
				"foo": &zcl.Attribute{
					Name: "foo",
					Expr: &expression{
						src: &hclast.ObjectType{
							List: &hclast.ObjectList{},
							Lbrace: hcltoken.Pos{
								Offset: 4,
								Line:   1,
								Column: 5,
							},
							Rbrace: hcltoken.Pos{
								Offset: 5,
								Line:   1,
								Column: 6,
							},
						},
					},
					Range: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
					NameRange: zcl.Range{
						Start: zcl.Pos{Byte: 0, Line: 1, Column: 1},
						End:   zcl.Pos{Byte: 1, Line: 1, Column: 2},
					},
				},
			},
			1, // warning about using block syntax
		},
		{
			`foo "bar" {}`,
			zcl.Attributes{},
			1, // blocks are not allowed here
		},
		{
			`
			    foo = 1
			    foo = 2
			`,
			zcl.Attributes{
				"foo": &zcl.Attribute{
					Name: "foo",
					Expr: &expression{
						src: &hclast.LiteralType{
							Token: hcltoken.Token{
								Type: hcltoken.NUMBER,
								Pos: hcltoken.Pos{
									Offset: 14,
									Line:   2,
									Column: 14,
								},
								Text: `1`,
							},
						},
					},
					Range: zcl.Range{
						Start: zcl.Pos{Byte: 8, Line: 2, Column: 8},
						End:   zcl.Pos{Byte: 9, Line: 2, Column: 9},
					},
					NameRange: zcl.Range{
						Start: zcl.Pos{Byte: 8, Line: 2, Column: 8},
						End:   zcl.Pos{Byte: 9, Line: 2, Column: 9},
					},
				},
			},
			1, // duplicate definition of foo
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			file, diags := Parse([]byte(test.Source), "test.hcl")
			if len(diags) != 0 {
				t.Fatalf("diagnostics from parse: %s", diags.Error())
			}

			got, diags := file.Body.JustAttributes()
			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf(" - %s", diag.Error())
				}
			}

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("wrong result\ngot:  %s\nwant: %s", spew.Sdump(got), spew.Sdump(test.Want))
			}
		})
	}
}

func TestExpressionValue(t *testing.T) {
	tests := []struct {
		Source    string // HCL source assigning a value to attribute "v"
		Want      cty.Value
		DiagCount int
	}{
		{
			`v = 1`,
			cty.NumberIntVal(1),
			0,
		},
		{
			`v = 1.5`,
			cty.NumberFloatVal(1.5),
			0,
		},
		{
			`v = "hello"`,
			cty.StringVal("hello"),
			0,
		},
		{
			`v = <<EOT
heredoc
EOT
`,
			cty.StringVal("heredoc\n"),
			0,
		},
		{
			`v = true`,
			cty.True,
			0,
		},
		{
			`v = false`,
			cty.False,
			0,
		},
		{
			`v = []`,
			cty.EmptyTupleVal,
			0,
		},
		{
			`v = ["hello", 5, true, 3.4]`,
			cty.TupleVal([]cty.Value{
				cty.StringVal("hello"),
				cty.NumberIntVal(5),
				cty.True,
				cty.NumberFloatVal(3.4),
			}),
			0,
		},
		{
			`v = {}`,
			cty.EmptyObjectVal,
			0,
		},
		{
			`v = {
				string = "hello"
				int = 5
				bool = true
				float = 3.4
				list = []
				object = {}
			}`,
			cty.ObjectVal(map[string]cty.Value{
				"string": cty.StringVal("hello"),
				"int":    cty.NumberIntVal(5),
				"bool":   cty.True,
				"float":  cty.NumberFloatVal(3.4),
				"list":   cty.EmptyTupleVal,
				"object": cty.EmptyObjectVal,
			}),
			0,
		},
		{
			`v {}`,
			cty.EmptyObjectVal,
			0, // warns about using block syntax during content extraction, but we ignore that here
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			file, diags := Parse([]byte(test.Source), "test.hcl")
			if len(diags) != 0 {
				t.Fatalf("diagnostics from parse: %s", diags.Error())
			}

			content, diags := file.Body.Content(&zcl.BodySchema{
				Attributes: []zcl.AttributeSchema{
					{
						Name:     "v",
						Required: true,
					},
				},
			})

			expr := content.Attributes["v"].Expr

			got, diags := expr.Value(nil)
			if len(diags) != test.DiagCount {
				t.Errorf("wrong number of diagnostics %d; want %d", len(diags), test.DiagCount)
				for _, diag := range diags {
					t.Logf(" - %s", diag.Error())
				}
			}

			if !got.RawEquals(test.Want) {
				t.Errorf("wrong result\ngot:  %#v\nwant: %#v", got, test.Want)
			}
		})
	}

}
