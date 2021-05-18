package lex

type keyword string

const (
andKeyword keyword = "and"
orKeyword keyword = "or"
)

//selectKeyword keyword = "select"
//fromKeyword   keyword = "from"
//asKeyword     keyword = "as"
//tableKeyword  keyword = "table"
//createKeyword keyword = "create"
//insertKeyword keyword = "insert"
//intoKeyword   keyword = "into"
//valuesKeyword keyword = "values"
//intKeyword    keyword = "int"
//textKeyword   keyword = "text"

type symbol string

const (
	semicolonSymbol  symbol = ";"
	asteriskSymbol   symbol = "*"
	commaSymbol      symbol = ","
	leftparenSymbol  symbol = "("
	rightparenSymbol symbol = ")"
)

