package common

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"testing"

	"github.com/xuri/excelize/v2"
)

// 定义一个泛型函数 excelFile，用于生成 Excel 文件
// T 是任意类型，表示每一行的数据类型
// 参数：
// - sheet: 工作表名称
// - columns: 列标题数组
// - rowValue: 每一行的数据数组
// - filename: 文件名
func excelFile[T any](sheet string, columns []string, rowValue []T, filename string) (string, error) {
	// 创建一个新的 Excel 文件
	f := excelize.NewFile()

	// 新建一个工作表，并获取其索引
	// index := f.NewSheet(sheet)
	index := f.GetSheetIndex(sheet) // 打开 sheet1

	// 遍历列标题数组，设置每列的标题
	for i := 0; i < len(columns); i++ {
		// 计算列标题的单元格坐标，例如 A1, B1, C1 等
		axis := string(letter[i]) + strconv.Itoa(1)
		// 在指定单元格设置列标题
		f.SetCellValue(sheet, axis, columns[i])
	}

	// 初始化行号，从第 2 行开始写入数据
	num := 2

	// 遍历每一行的数据
	for r := 0; r < len(rowValue); r++ {
		// 获取当前行的数据值
		v := reflect.ValueOf(rowValue[r])

		// 遍历当前行的每个字段
		c := v.NumField()
		fmt.Println(c)
		for i := 0; i < v.NumField(); i++ {
			// 计算当前字段的单元格坐标，例如 A2, B2, C2 等
			axis := string(letter[i]) + strconv.Itoa(num)
			// 在指定单元格设置字段值
			f.SetCellValue(sheet, axis, v.Field(i).Interface())
		}

		// 行号递增，准备写入下一行数据
		num++
	}

	// 设置活动工作表为新创建的工作表
	f.SetActiveSheet(index)

	// 构建保存文件的路径
	path := "/Users/zp/zptmp/tmp/" + filename + ".xlsx"

	// 将 Excel 文件保存到指定路径
	if err := f.SaveAs(path); err != nil {
		// 如果保存文件时出错，记录错误日志并返回错误
		log.Printf("[ERROR] file save error: %s", err)
		return "", err
	}

	// 返回保存文件的路径
	return path, nil
}

func TestExcel(t *testing.T) {

	sheet := "sheet1"
	columns := []string{"合同名称", "客户名称", "合同金额"}
	fileName := "test_export"

	type Info struct {
		ContractName  string  `json:"合同名称"`
		CustomerName  string  `json:"客户名称"`
		ContractMoney float64 `json:"合同金额"`
	}

	data := []Info{
		{"合同1", "客户1", 100},
		{"合同2", "客户2", 200},
		{"合同3", "客户3", 300},
	}

	file, err := excelFile(sheet, columns, data, fileName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(file)
}
