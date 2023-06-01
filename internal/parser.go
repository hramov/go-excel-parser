package internal

import (
	"encoding/json"
	"fmt"
	"github.com/hramov/go-excel-parser/pkg/utils"
	"github.com/xuri/excelize/v2"
	"log"
)

type parser struct {
	Path     string
	Buffer   []byte
	Template []byte
}

func NewFileParser(path string, template []byte) Parser {
	return &parser{
		Path:     path,
		Template: template,
	}
}

func (p *parser) Parse() ([]byte, error) {
	return p.parseExcelToJson()
}

func (p *parser) parseExcelToJson() ([]byte, error) {
	if p.Path == "" {
		return nil, fmt.Errorf("parser: no path provided")
	}

	var template TemplateConfig
	err := json.Unmarshal(p.Template, &template)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenFile(p.Path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}()

	var rows [][]string

	if template.SheetName == "" && template.SheetNumber != 0 {
		sheetName := f.GetSheetList()[template.SheetNumber-1]
		if sheetName == "" {
			return nil, fmt.Errorf("parser: cannot fetch sheet name by index")
		}
		template.SheetName = sheetName
	}

	if template.SheetName == "" {
		return nil, fmt.Errorf("parser: no sheet name provided or computed")
	}

	rows, err = f.GetRows(template.SheetName, excelize.Options{
		RawCellValue: true,
	})
	if err != nil {
		return nil, err
	}

	template.SetFieldsStr()

	var headers []string
	var result []map[string]string

	for index, row := range rows {
		if template.StartRow > index+1 {
			if index+1 == template.HeaderRow {
				for _, headerCell := range row {
					fieldIndex := utils.Includes(template.FieldsStr, headerCell)
					if fieldIndex > -1 {
						headers = append(headers, template.Fields[fieldIndex].Name)
					} else {
						return nil, fmt.Errorf("parser: no such column %s in config", headerCell)
					}
				}
			}
			continue
		}

		m := make(map[string]string)
		for dataIndex, dataCell := range row {
			m[headers[dataIndex]] = dataCell
		}

		result = append(result, m)
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return jsonBytes, nil
}
