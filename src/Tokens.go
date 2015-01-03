package main

type Token int

const (
	TOKEN_EQ          Token = 1 << iota // =
	TOKEN_EQEQ                          // ==
	TOKEN_LT                            // <
	TOKEN_GT                            // >
	TOKEN_LEFT_PAREN                    // (
	TOKEN_RIGHT_PAREN                   // )
	TOKEN_LEFT_BRACE                    // {
	TOKEN_RIGHT_BRACE                   // }
	TOKEN_NEWLINE                       // \n
	TOKEN_VAR                           // var

	TOKEN_IF
	TOKEN_ELSE

	TOKEN_PLUS
	TOKEN_MINUS

	TOKEN_PLUSEQ  // abc += 3
	TOKEN_MINUSEQ // abc -= 3

	TOKEN_VALUE // 123
	TOKEN_NAME  // abc

	TOKEN_LEXER_ERROR // Used when the Lexer made a booboo
	TOKEN_END_OF_PROGRAM
)

type GusToken struct {
	Token Token
	Value string
}
