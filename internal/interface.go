package internal

const (
	Text   = "text"
	Number = "number"
	Date   = "date"
)

type Parser interface {
	Parse() ([]byte, error)
}

type TemplateConfig struct {
	Id          int8    `json:"id"`
	Title       string  `json:"title"`
	SheetName   string  `json:"sheet_name"`
	SheetNumber int     `json:"sheet_number"`
	HeaderRow   int     `json:"header_row"`
	StartRow    int     `json:"start_row"`
	Fields      []Field `json:"fields"`
	FieldsStr   []string
}

type Field struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Field string `json:"field"`
	Type  string `json:"type"`
}

func (t *TemplateConfig) SetFieldsStr() {
	for _, field := range t.Fields {
		t.FieldsStr = append(t.FieldsStr, field.Field)
	}
}
