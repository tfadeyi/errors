package grammar

import "github.com/alecthomas/participle/v2/lexer"

var lexerDefinition = lexer.MustSimple([]lexer.SimpleRule{
	{"EOL", `[\n\r]+`},
	{"Fyi", `@fyi`},
	{"String", `([a-zA-Z_0-9\.\/:,\-\'\(\)~\[\]\{\}=\"\|%])\w*`},
	{"Whitespace", `[ \t]+`},
})
