package resmodels

type Currency struct {
	Code            string  `json:"code"`
	CodeNumeric     int     `json:"codeNumeric"`
	Digits          int     `json:"digits"`
	EuroCR          float64 `json:"euroCr"`
	EuroCRUpdatedAt string  `json:"euroCrUpdatedAt"`
	Symbol          string  `json:"symbol"`
}
