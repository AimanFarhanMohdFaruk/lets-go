package ui

type InputConfig struct {
	Typ string "json:\"type\""
	Name string "json:\"name\""
	Label string "json:\"label\""
	Placeholder string "json:\"placeholder\""
	Required bool "json:\"required\""
	Err string "json:\"error\""
	Value string "json:\"value\""
	Checked bool
}