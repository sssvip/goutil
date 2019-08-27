package excelutil

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sssvip/goutil/logutil"
	"github.com/sssvip/goutil/strutil"
	"strconv"
)

type ExcelWrapper struct {
	name         string
	excel        *excelize.File
	sheetLineMap map[string]int
}

func NewExcel(fileName string) *ExcelWrapper {
	return &ExcelWrapper{fileName, excelize.NewFile(), make(map[string]int)}
}

func (wrapper *ExcelWrapper) SetTitle() {

}
func (wrapper *ExcelWrapper) GetExcelOrigin() *excelize.File {
	return wrapper.excel
}

func (wrapper *ExcelWrapper) AppendLine(strs ...string) {
	wrapper.AppendLineForSheet("sheet1", strs...)
}
func (wrapper *ExcelWrapper) AppendLineForSheet(sheet string, strs ...string) {
	currentLine, ok := wrapper.sheetLineMap[sheet]
	if !ok {
		currentLine = 1
	}
	for i, s := range strs {
		cellIndex := string(rune(65+i)) + strconv.Itoa(currentLine)
		wrapper.excel.SetCellStr(sheet, cellIndex, s)
	}
	currentLine++
	wrapper.sheetLineMap[sheet] = currentLine
}

func (wrapper *ExcelWrapper) Save() {
	err := wrapper.excel.SaveAs(strutil.Format("%s.xlsx", wrapper.name))
	if err != nil {
		logutil.Error.Println(err)
	}
}
