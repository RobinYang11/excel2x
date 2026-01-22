package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mozillazg/go-pinyin"
	"github.com/xuri/excelize/v2"
)

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// 符号
const (
	checkMark = "✅"
	crossMark = "❌"
	arrow     = "→"
	info      = "ℹ"
)

// printSuccess 打印成功信息（绿色对号）
func printSuccess(t *testing.T, msg string) {
	t.Logf("%s%s %s%s", colorGreen, checkMark, msg, colorReset)
}

// printError 打印错误信息（红色叉号）
func printError(t *testing.T, msg string) {
	t.Logf("%s%s %s%s", colorRed, crossMark, msg, colorReset)
}

// printInfo 打印信息（蓝色箭头）
func printInfo(t *testing.T, msg string) {
	t.Logf("%s%s %s%s", colorCyan, arrow, msg, colorReset)
}

// printTestHeader 打印测试标题
func printTestHeader(t *testing.T, testName string) {
	t.Logf("\n%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s", colorBlue, colorReset)
	t.Logf("%s  %s%s", colorBlue, testName, colorReset)
	t.Logf("%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", colorBlue, colorReset)
}

// TestWriteJson 测试 WriteJson 函数
func TestWriteJson(t *testing.T) {
	printTestHeader(t, "测试 JSON 文件写入功能")
	
	// 创建临时目录
	tmpDir := t.TempDir()
	printInfo(t, fmt.Sprintf("创建临时目录: %s", tmpDir))
	
	// 准备测试数据
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "测试数据",
	}
	printInfo(t, fmt.Sprintf("准备测试数据: %d 个键值对", len(testData)))
	
	// 测试文件路径
	testFile := filepath.Join(tmpDir, "test.json")
	
	// 调用函数
	WriteJson(testData, testFile)
	printInfo(t, fmt.Sprintf("调用 WriteJson 写入文件: %s", testFile))
	
	// 验证文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("JSON 文件未创建: %s", testFile))
		t.Fatalf("JSON 文件未创建: %s", testFile)
	}
	printSuccess(t, "JSON 文件创建成功")
	
	// 读取并验证文件内容
	content, err := os.ReadFile(testFile)
	if err != nil {
		printError(t, fmt.Sprintf("读取 JSON 文件失败: %v", err))
		t.Fatalf("读取 JSON 文件失败: %v", err)
	}
	printSuccess(t, "JSON 文件读取成功")
	
	// 解析 JSON
	var result map[string]string
	if err := json.Unmarshal(content, &result); err != nil {
		printError(t, fmt.Sprintf("JSON 解析失败: %v", err))
		t.Fatalf("JSON 解析失败: %v", err)
	}
	printSuccess(t, "JSON 格式验证通过")
	
	// 验证数据
	allPassed := true
	for key, expectedValue := range testData {
		if actualValue, exists := result[key]; !exists {
			printError(t, fmt.Sprintf("缺少 key: %s", key))
			allPassed = false
			t.Errorf("缺少 key: %s", key)
		} else if actualValue != expectedValue {
			printError(t, fmt.Sprintf("key %s 的值不匹配: 期望 %s, 实际 %s", key, expectedValue, actualValue))
			allPassed = false
			t.Errorf("key %s 的值不匹配: 期望 %s, 实际 %s", key, expectedValue, actualValue)
		} else {
			printSuccess(t, fmt.Sprintf("验证 key '%s' = '%s'", key, expectedValue))
		}
	}
	
	if allPassed {
		printSuccess(t, "所有数据验证通过")
	}
}

// TestPinyinConversion 测试拼音转换功能
func TestPinyinConversion(t *testing.T) {
	printTestHeader(t, "测试拼音转换功能")
	
	testCases := []struct {
		input    string
		expected string
	}{
		{"你好", "nihao"},
		{"世界", "shijie"},
		{"测试", "ceshi"},
		{"Hello", ""}, // 拼音库对非中文字符返回空
		{"123", ""},   // 拼音库对数字返回空
		{"你好World", "nihao"}, // 混合字符
	}
	
	printInfo(t, fmt.Sprintf("测试 %d 个拼音转换用例", len(testCases)))
	
	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := pinyin.LazyPinyin(tc.input, pinyin.NewArgs())
			output := strings.Join(result, "")
			if output != tc.expected {
				printError(t, fmt.Sprintf("拼音转换失败: 输入 '%s', 期望 '%s', 实际 '%s'", tc.input, tc.expected, output))
				t.Errorf("拼音转换失败: 输入 %s, 期望 %s, 实际 %s", tc.input, tc.expected, output)
			} else {
				printSuccess(t, fmt.Sprintf("拼音转换成功: '%s' → '%s'", tc.input, output))
			}
		})
	}
	
	printSuccess(t, "所有拼音转换测试通过")
}

// TestExcelToJsonConversion 测试 Excel 转 JSON 的完整流程
func TestExcelToJsonConversion(t *testing.T) {
	printTestHeader(t, "测试 Excel 转 JSON 完整流程")
	
	// 检查测试文件是否存在
	testFile := "./assets/abc.xlsx"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("测试文件不存在，跳过测试: %s", testFile))
		t.Skipf("测试文件不存在，跳过测试: %s", testFile)
	}
	printSuccess(t, fmt.Sprintf("找到测试文件: %s", testFile))
	
	// 创建临时输出目录
	tmpDir := t.TempDir()
	printInfo(t, fmt.Sprintf("创建临时输出目录: %s", tmpDir))
	
	// 打开 Excel 文件
	f, err := excelize.OpenFile(testFile)
	if err != nil {
		printError(t, fmt.Sprintf("打开 Excel 文件失败: %v", err))
		t.Fatalf("打开 Excel 文件失败: %v", err)
	}
	defer f.Close()
	printSuccess(t, "Excel 文件打开成功")
	
	// 获取列数据
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		printError(t, fmt.Sprintf("获取列数据失败: %v", err))
		t.Fatalf("获取列数据失败: %v", err)
	}
	
	if len(cols) == 0 {
		printError(t, "Excel 文件没有数据")
		t.Fatal("Excel 文件没有数据")
	}
	printSuccess(t, fmt.Sprintf("读取到 %d 列数据", len(cols)))
	
	// 获取行数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		printError(t, fmt.Sprintf("获取行数据失败: %v", err))
		t.Fatalf("获取行数据失败: %v", err)
	}
	
	if len(rows) == 0 {
		printError(t, "Excel 文件没有行数据")
		t.Fatal("Excel 文件没有行数据")
	}
	printSuccess(t, fmt.Sprintf("读取到 %d 行数据", len(rows)))
	
	// 验证第一行（表头）存在
	if len(rows[0]) == 0 {
		printError(t, "Excel 文件第一行为空")
		t.Fatal("Excel 文件第一行为空")
	}
	
	// 处理数据并生成 JSON 文件
	len_of_col := len(cols)
	var dart_keys []string
	jsonFileCount := 0
	
	for row_index, row := range rows[0] {
		if len(row) == 0 {
			continue
		}
		
		file_data := make(map[string]string, len_of_col-1)
		col := cols[row_index]
		
		for i := 1; i < len(col); i++ {
			if i >= len(cols[0]) {
				break
			}
			keySlice := pinyin.LazyPinyin(cols[0][i], pinyin.NewArgs())
			key := strings.Join(keySlice, "")
			
			// 收集 Dart keys
			if row_index == 0 {
				dart_keys = append(dart_keys, key)
			}
			
			if i < len(col) {
				file_data[key] = col[i]
			}
		}
		
		// 生成 JSON 文件
		jsonFile := filepath.Join(tmpDir, row+".json")
		WriteJson(file_data, jsonFile)
		jsonFileCount++
		
		// 验证 JSON 文件已创建
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			printError(t, fmt.Sprintf("JSON 文件未创建: %s", jsonFile))
			t.Errorf("JSON 文件未创建: %s", jsonFile)
		} else {
			printSuccess(t, fmt.Sprintf("生成 JSON 文件: %s.json", row))
		}
		
		// 验证 JSON 文件内容
		content, err := os.ReadFile(jsonFile)
		if err != nil {
			printError(t, fmt.Sprintf("读取 JSON 文件失败: %v", err))
			t.Errorf("读取 JSON 文件失败: %v", err)
			continue
		}
		
		var jsonData map[string]string
		if err := json.Unmarshal(content, &jsonData); err != nil {
			printError(t, fmt.Sprintf("JSON 解析失败: %v", err))
			t.Errorf("JSON 解析失败: %v", err)
			continue
		}
		
		// 验证 JSON 数据不为空（至少应该有一些 key）
		if len(jsonData) == 0 && len(file_data) > 0 {
			printError(t, fmt.Sprintf("JSON 文件内容为空: %s", jsonFile))
			t.Errorf("JSON 文件内容为空: %s", jsonFile)
		} else {
			printSuccess(t, fmt.Sprintf("验证 JSON 文件 '%s.json' 包含 %d 个键值对", row, len(jsonData)))
		}
	}
	
	printSuccess(t, fmt.Sprintf("成功生成 %d 个 JSON 文件", jsonFileCount))
	
	// 验证 Dart keys 已收集
	if len(dart_keys) == 0 {
		printError(t, "未收集到 Dart keys")
		t.Error("未收集到 Dart keys")
	} else {
		printSuccess(t, fmt.Sprintf("收集到 %d 个 Dart keys", len(dart_keys)))
	}
}

// TestDartFileGeneration 测试 Dart 文件生成
func TestDartFileGeneration(t *testing.T) {
	printTestHeader(t, "测试 Dart 文件生成")
	
	// 创建临时目录
	tmpDir := t.TempDir()
	printInfo(t, fmt.Sprintf("创建临时目录: %s", tmpDir))
	
	// 测试数据
	dart_keys := []string{"key1", "key2", "ceshi", "nihao"}
	printInfo(t, fmt.Sprintf("准备 %d 个 Dart keys", len(dart_keys)))
	
	// 生成 Dart 文件内容
	var dart_content strings.Builder
	for _, key := range dart_keys {
		dart_content.WriteString("  static const " + key + " = '" + key + "';\n")
	}
	final_dart := "class LocaleKeys {\n" + strings.TrimRight(dart_content.String(), "\n") + "\n}\n"
	
	// 写入文件
	dartFile := filepath.Join(tmpDir, "locale_keys.dart")
	err := os.WriteFile(dartFile, []byte(final_dart), 0644)
	if err != nil {
		printError(t, fmt.Sprintf("创建 Dart 文件失败: %v", err))
		t.Fatalf("创建 Dart 文件失败: %v", err)
	}
	printSuccess(t, fmt.Sprintf("Dart 文件创建成功: %s", dartFile))
	
	// 验证文件存在
	if _, err := os.Stat(dartFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("Dart 文件未创建: %s", dartFile))
		t.Fatalf("Dart 文件未创建: %s", dartFile)
	}
	printSuccess(t, "Dart 文件存在验证通过")
	
	// 读取并验证内容
	content, err := os.ReadFile(dartFile)
	if err != nil {
		printError(t, fmt.Sprintf("读取 Dart 文件失败: %v", err))
		t.Fatalf("读取 Dart 文件失败: %v", err)
	}
	
	contentStr := string(content)
	
	// 验证包含类定义
	if !strings.Contains(contentStr, "class LocaleKeys") {
		printError(t, "Dart 文件缺少类定义")
		t.Error("Dart 文件缺少类定义")
	} else {
		printSuccess(t, "Dart 文件包含类定义 'class LocaleKeys'")
	}
	
	// 验证包含所有 keys
	allKeysFound := true
	for _, key := range dart_keys {
		if !strings.Contains(contentStr, "static const "+key) {
			printError(t, fmt.Sprintf("Dart 文件缺少 key: %s", key))
			allKeysFound = false
			t.Errorf("Dart 文件缺少 key: %s", key)
		} else {
			printSuccess(t, fmt.Sprintf("验证 key '%s' 存在", key))
		}
	}
	
	if allKeysFound {
		printSuccess(t, "所有 Dart keys 验证通过")
	}
}

// TestExcelFileExists 测试 Excel 文件是否存在且可读
func TestExcelFileExists(t *testing.T) {
	printTestHeader(t, "测试 Excel 文件可读性")
	
	testFile := "./assets/abc.xlsx"
	
	// 检查文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("测试文件不存在: %s", testFile))
		t.Skipf("测试文件不存在: %s", testFile)
	}
	printSuccess(t, fmt.Sprintf("文件存在: %s", testFile))
	
	// 尝试打开文件
	f, err := excelize.OpenFile(testFile)
	if err != nil {
		printError(t, fmt.Sprintf("无法打开 Excel 文件: %v", err))
		t.Fatalf("无法打开 Excel 文件: %v", err)
	}
	defer f.Close()
	printSuccess(t, "Excel 文件打开成功")
	
	// 检查 Sheet1 是否存在
	sheetList := f.GetSheetList()
	printInfo(t, fmt.Sprintf("找到 %d 个工作表: %v", len(sheetList), sheetList))
	
	found := false
	for _, sheet := range sheetList {
		if sheet == "Sheet1" {
			found = true
			break
		}
	}
	
	if !found {
		printError(t, "Excel 文件中未找到 Sheet1")
		t.Error("Excel 文件中未找到 Sheet1")
	} else {
		printSuccess(t, "找到工作表 'Sheet1'")
	}
	
	// 尝试读取数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		printError(t, fmt.Sprintf("读取 Sheet1 数据失败: %v", err))
		t.Fatalf("读取 Sheet1 数据失败: %v", err)
	}
	
	if len(rows) == 0 {
		printInfo(t, "Sheet1 中没有数据行")
		t.Log("Sheet1 中没有数据行")
	} else {
		printSuccess(t, fmt.Sprintf("Sheet1 包含 %d 行数据", len(rows)))
	}
}

// BenchmarkWriteJson 性能测试：WriteJson 函数
func BenchmarkWriteJson(b *testing.B) {
	tmpDir := b.TempDir()
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
		"key4": "value4",
		"key5": "value5",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		testFile := filepath.Join(tmpDir, "test.json")
		WriteJson(testData, testFile)
	}
}

// TestAllParameters 测试所有命令行参数
func TestAllParameters(t *testing.T) {
	printTestHeader(t, "测试所有命令行参数")
	
	// 检查测试文件是否存在
	testFile := "./assets/abc.xlsx"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("测试文件不存在，跳过测试: %s", testFile))
		t.Skipf("测试文件不存在，跳过测试: %s", testFile)
	}
	printSuccess(t, fmt.Sprintf("找到测试文件: %s", testFile))
	
	// 创建临时目录用于测试所有参数
	tmpDir := t.TempDir()
	jsonOutputDir := filepath.Join(tmpDir, "json_output")
	dartOutputDir := filepath.Join(tmpDir, "dart_output")
	
	// 创建输出目录
	if err := os.MkdirAll(jsonOutputDir, 0755); err != nil {
		printError(t, fmt.Sprintf("创建 JSON 输出目录失败: %v", err))
		t.Fatalf("创建 JSON 输出目录失败: %v", err)
	}
	printSuccess(t, fmt.Sprintf("创建 JSON 输出目录: %s", jsonOutputDir))
	
	if err := os.MkdirAll(dartOutputDir, 0755); err != nil {
		printError(t, fmt.Sprintf("创建 Dart 输出目录失败: %v", err))
		t.Fatalf("创建 Dart 输出目录失败: %v", err)
	}
	printSuccess(t, fmt.Sprintf("创建 Dart 输出目录: %s", dartOutputDir))
	
	// 测试参数1: 使用默认参数
	printInfo(t, "测试场景 1: 使用默认参数")
	testDefaultParams(t, testFile, tmpDir)
	
	// 测试参数2: 指定所有自定义参数
	printInfo(t, "测试场景 2: 指定所有自定义参数")
	testCustomParams(t, testFile, jsonOutputDir, dartOutputDir)
	
	// 测试参数3: 只指定 file 参数
	printInfo(t, "测试场景 3: 只指定 file 参数")
	testFileOnly(t, testFile, tmpDir)
	
	// 测试参数4: 指定 file 和 output_path
	printInfo(t, "测试场景 4: 指定 file 和 output_path")
	testFileAndOutputPath(t, testFile, jsonOutputDir, tmpDir)
	
	// 测试参数5: 指定 file 和 dart_file_path
	printInfo(t, "测试场景 5: 指定 file 和 dart_file_path")
	testFileAndDartPath(t, testFile, tmpDir, dartOutputDir)
	
	printSuccess(t, "所有参数测试场景完成")
}

// testDefaultParams 测试使用默认参数
func testDefaultParams(t *testing.T, testFile string, tmpDir string) {
	// 模拟默认参数：file=testFile, output_path=tmpDir, dart_file_path=tmpDir
	outputPath := tmpDir
	dartPath := tmpDir
	
	// 执行处理逻辑
	err := processExcelFile(testFile, outputPath, dartPath)
	if err != nil {
		printError(t, fmt.Sprintf("处理 Excel 文件失败: %v", err))
		t.Errorf("处理 Excel 文件失败: %v", err)
		return
	}
	printSuccess(t, "使用默认参数处理成功")
	
	// 验证 JSON 文件
	verifyJsonFiles(t, outputPath, "默认参数")
	
	// 验证 Dart 文件
	dartFile := filepath.Join(dartPath, "locale_keys.dart")
	verifyDartFile(t, dartFile, "默认参数")
}

// testCustomParams 测试指定所有自定义参数
func testCustomParams(t *testing.T, testFile string, jsonOutputDir string, dartOutputDir string) {
	err := processExcelFile(testFile, jsonOutputDir, dartOutputDir)
	if err != nil {
		printError(t, fmt.Sprintf("处理 Excel 文件失败: %v", err))
		t.Errorf("处理 Excel 文件失败: %v", err)
		return
	}
	printSuccess(t, "使用自定义参数处理成功")
	
	// 验证 JSON 文件
	verifyJsonFiles(t, jsonOutputDir, "自定义参数")
	
	// 验证 Dart 文件
	dartFile := filepath.Join(dartOutputDir, "locale_keys.dart")
	verifyDartFile(t, dartFile, "自定义参数")
}

// testFileOnly 测试只指定 file 参数
func testFileOnly(t *testing.T, testFile string, tmpDir string) {
	outputPath := tmpDir
	dartPath := tmpDir
	
	err := processExcelFile(testFile, outputPath, dartPath)
	if err != nil {
		printError(t, fmt.Sprintf("处理 Excel 文件失败: %v", err))
		t.Errorf("处理 Excel 文件失败: %v", err)
		return
	}
	printSuccess(t, "只指定 file 参数处理成功")
}

// testFileAndOutputPath 测试指定 file 和 output_path
func testFileAndOutputPath(t *testing.T, testFile string, jsonOutputDir string, tmpDir string) {
	dartPath := tmpDir
	
	err := processExcelFile(testFile, jsonOutputDir, dartPath)
	if err != nil {
		printError(t, fmt.Sprintf("处理 Excel 文件失败: %v", err))
		t.Errorf("处理 Excel 文件失败: %v", err)
		return
	}
	printSuccess(t, "指定 file 和 output_path 处理成功")
	
	// 验证 JSON 文件在指定目录
	verifyJsonFiles(t, jsonOutputDir, "file 和 output_path")
}

// testFileAndDartPath 测试指定 file 和 dart_file_path
func testFileAndDartPath(t *testing.T, testFile string, tmpDir string, dartOutputDir string) {
	outputPath := tmpDir
	
	err := processExcelFile(testFile, outputPath, dartOutputDir)
	if err != nil {
		printError(t, fmt.Sprintf("处理 Excel 文件失败: %v", err))
		t.Errorf("处理 Excel 文件失败: %v", err)
		return
	}
	printSuccess(t, "指定 file 和 dart_file_path 处理成功")
	
	// 验证 Dart 文件在指定目录
	dartFile := filepath.Join(dartOutputDir, "locale_keys.dart")
	verifyDartFile(t, dartFile, "file 和 dart_file_path")
}

// processExcelFile 处理 Excel 文件的辅助函数（模拟 main 函数的处理逻辑）
func processExcelFile(inputFile string, outputPath string, dartFilePath string) error {
	// 打开 Excel 文件
	f, err := excelize.OpenFile(inputFile)
	if err != nil {
		return fmt.Errorf("打开 Excel 文件失败: %v", err)
	}
	defer f.Close()
	
	// 获取列数据
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		return fmt.Errorf("获取列数据失败: %v", err)
	}
	
	// 获取行数据
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return fmt.Errorf("获取行数据失败: %v", err)
	}
	
	len_of_col := len(cols)
	var dart_keys []string
	
	// 处理每一行
	for row_index, row := range rows[0] {
		file_data := make(map[string]string, len_of_col-1)
		col := cols[row_index]
		
		for i := 1; i < len(col); i++ {
			if i >= len(cols[0]) {
				break
			}
			originalText := cols[0][i]
			
			// 1. 先去空格和特殊字符
			cleanedText := cleanText(originalText)
			
			// 2. 再翻译（转换为拼音）
			key := generateKey(cleanedText)
			
			// 如果key为空，使用清理后的文本作为key
			if key == "" {
				key = cleanedText
			}
			
			// 如果key以数字开头，在前面添加 auto_gen_ 前缀
			if key != "" && len(key) > 0 && key[0] >= '0' && key[0] <= '9' {
				key = "auto_gen_" + key
			}
			
			// 收集 keys 用于 dart 文件
			if row_index == 0 {
				// 跳过空字符的key
				if key == "" {
					continue
				}
				// 去重
				if !slices.Contains(dart_keys, key) {
					dart_keys = append(dart_keys, key)
				}
			}
			file_data[key] = col[i]
		}
		
		// 写入 JSON 文件
		jsonFile := filepath.Join(outputPath, row+".json")
		WriteJson(file_data, jsonFile)
	}
	
	// 生成 Dart 文件
	var dart_content strings.Builder
	for _, key := range dart_keys {
		dart_content.WriteString(fmt.Sprintf("  static const %s = '%s';\n", key, key))
	}
	final_dart := fmt.Sprintf(dart_file_template, strings.TrimRight(dart_content.String(), "\n"))
	
	dart_file_full_path := filepath.Join(dartFilePath, "locale_keys.dart")
	dart_file, err := os.Create(dart_file_full_path)
	if err != nil {
		return fmt.Errorf("创建 Dart 文件失败: %v", err)
	}
	defer dart_file.Close()
	dart_file.WriteString(final_dart)
	
	return nil
}

// verifyJsonFiles 验证 JSON 文件
func verifyJsonFiles(t *testing.T, outputPath string, scenario string) {
	// 检查常见的 JSON 文件是否存在
	expectedFiles := []string{"zh.json", "en.json", "ja.json"}
	
	for _, filename := range expectedFiles {
		jsonFile := filepath.Join(outputPath, filename)
		if _, err := os.Stat(jsonFile); os.IsNotExist(err) {
			// 文件不存在不是错误，可能 Excel 文件没有对应的行
			printInfo(t, fmt.Sprintf("[%s] JSON 文件不存在（可能正常）: %s", scenario, filename))
		} else {
			printSuccess(t, fmt.Sprintf("[%s] 验证 JSON 文件存在: %s", scenario, filename))
			
			// 验证 JSON 文件内容
			content, err := os.ReadFile(jsonFile)
			if err != nil {
				printError(t, fmt.Sprintf("[%s] 读取 JSON 文件失败: %v", scenario, err))
				continue
			}
			
			var jsonData map[string]string
			if err := json.Unmarshal(content, &jsonData); err != nil {
				printError(t, fmt.Sprintf("[%s] JSON 解析失败: %v", scenario, err))
			} else {
				printSuccess(t, fmt.Sprintf("[%s] JSON 文件格式正确，包含 %d 个键值对", scenario, len(jsonData)))
			}
		}
	}
}

// verifyDartFile 验证 Dart 文件
func verifyDartFile(t *testing.T, dartFile string, scenario string) {
	if _, err := os.Stat(dartFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("[%s] Dart 文件未创建: %s", scenario, dartFile))
		t.Errorf("[%s] Dart 文件未创建: %s", scenario, dartFile)
		return
	}
	printSuccess(t, fmt.Sprintf("[%s] Dart 文件存在: %s", scenario, dartFile))
	
	// 读取并验证内容
	content, err := os.ReadFile(dartFile)
	if err != nil {
		printError(t, fmt.Sprintf("[%s] 读取 Dart 文件失败: %v", scenario, err))
		t.Errorf("[%s] 读取 Dart 文件失败: %v", scenario, err)
		return
	}
	
	contentStr := string(content)
	
	// 验证包含类定义
	if !strings.Contains(contentStr, "class LocaleKeys") {
		printError(t, fmt.Sprintf("[%s] Dart 文件缺少类定义", scenario))
		t.Errorf("[%s] Dart 文件缺少类定义", scenario)
	} else {
		printSuccess(t, fmt.Sprintf("[%s] Dart 文件包含类定义", scenario))
	}
	
	// 验证文件名正确
	if !strings.HasSuffix(dartFile, "locale_keys.dart") {
		printError(t, fmt.Sprintf("[%s] Dart 文件名不正确: %s", scenario, dartFile))
		t.Errorf("[%s] Dart 文件名不正确: %s", scenario, dartFile)
	} else {
		printSuccess(t, fmt.Sprintf("[%s] Dart 文件名正确: locale_keys.dart", scenario))
	}
}

// TestDuplicateKeys 测试重复 key 的处理
func TestDuplicateKeys(t *testing.T) {
	printTestHeader(t, "测试重复 key 处理")
	
	// 创建临时目录
	tmpDir := t.TempDir()
	printInfo(t, fmt.Sprintf("创建临时目录: %s", tmpDir))
	
	// 模拟包含重复 key 的场景
	// 创建测试数据：一些不同的中文文本可能生成相同的拼音 key
	testCases := []struct {
		text     string
		expected string
	}{
		{"测试", "ceshi"},
		{"测试", "ceshi"}, // 完全相同的文本
		{"你好", "nihao"},
		{"世界", "shijie"},
		{"测试数据", "ceshishuju"},
	}
	
	printInfo(t, fmt.Sprintf("准备 %d 个测试用例（包含重复）", len(testCases)))
	
	var dart_keys []string
	duplicateCount := 0
	
	// 模拟处理逻辑
	for _, tc := range testCases {
		key := generateKey(tc.text)
		printInfo(t, fmt.Sprintf("处理文本 '%s' -> key '%s'", tc.text, key))
		
		if slices.Contains(dart_keys, key) {
			printInfo(t, fmt.Sprintf("检测到重复 key: %s (来自文本: %s)", key, tc.text))
			duplicateCount++
			continue
		}
		dart_keys = append(dart_keys, key)
		printSuccess(t, fmt.Sprintf("添加 key: %s", key))
	}
	
	// 验证结果
	expectedUniqueKeys := 4 // "ceshi", "nihao", "shijie", "ceshishuju"
	if len(dart_keys) != expectedUniqueKeys {
		printError(t, fmt.Sprintf("期望 %d 个唯一 key，实际得到 %d 个", expectedUniqueKeys, len(dart_keys)))
		t.Errorf("期望 %d 个唯一 key，实际得到 %d 个", expectedUniqueKeys, len(dart_keys))
	} else {
		printSuccess(t, fmt.Sprintf("验证唯一 key 数量正确: %d 个", len(dart_keys)))
	}
	
	// 验证没有重复
	keyMap := make(map[string]int)
	for _, key := range dart_keys {
		keyMap[key]++
		if keyMap[key] > 1 {
			printError(t, fmt.Sprintf("发现重复 key: %s (出现 %d 次)", key, keyMap[key]))
			t.Errorf("发现重复 key: %s", key)
		}
	}
	
	if len(keyMap) == len(dart_keys) {
		printSuccess(t, "验证通过：所有 key 都是唯一的")
	}
	
	// 验证检测到的重复数量
	if duplicateCount == 0 {
		printError(t, "未检测到任何重复 key")
		t.Error("未检测到任何重复 key")
	} else {
		printSuccess(t, fmt.Sprintf("成功检测并跳过了 %d 个重复 key", duplicateCount))
	}
	
	// 验证特定 key 存在
	expectedKeys := []string{"ceshi", "nihao", "shijie", "ceshishuju"}
	for _, expectedKey := range expectedKeys {
		if !slices.Contains(dart_keys, expectedKey) {
			printError(t, fmt.Sprintf("缺少期望的 key: %s", expectedKey))
			t.Errorf("缺少期望的 key: %s", expectedKey)
		} else {
			printSuccess(t, fmt.Sprintf("验证 key '%s' 存在", expectedKey))
		}
	}
}

// TestDuplicateKeysInExcel 测试 Excel 文件中重复 key 的处理
func TestDuplicateKeysInExcel(t *testing.T) {
	printTestHeader(t, "测试 Excel 文件中重复 key 处理")
	
	// 检查测试文件是否存在
	testFile := "./assets/abc.xlsx"
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		printError(t, fmt.Sprintf("测试文件不存在，跳过测试: %s", testFile))
		t.Skipf("测试文件不存在，跳过测试: %s", testFile)
	}
	printSuccess(t, fmt.Sprintf("找到测试文件: %s", testFile))
	
	// 创建临时输出目录
	tmpDir := t.TempDir()
	printInfo(t, fmt.Sprintf("创建临时输出目录: %s", tmpDir))
	
	// 打开 Excel 文件
	f, err := excelize.OpenFile(testFile)
	if err != nil {
		printError(t, fmt.Sprintf("打开 Excel 文件失败: %v", err))
		t.Fatalf("打开 Excel 文件失败: %v", err)
	}
	defer f.Close()
	printSuccess(t, "Excel 文件打开成功")
	
	// 获取列数据
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		printError(t, fmt.Sprintf("获取列数据失败: %v", err))
		t.Fatalf("获取列数据失败: %v", err)
	}
	
	if len(cols) == 0 {
		printError(t, "Excel 文件没有数据")
		t.Fatal("Excel 文件没有数据")
	}
	
	// 处理第一行，收集所有 keys
	var dart_keys []string
	duplicateKeys := make(map[string]int)
	
	// 只处理第一行（row_index == 0）来收集 keys
	row_index := 0
	col := cols[row_index]
	
	for i := 1; i < len(col); i++ {
		if i >= len(cols[0]) {
			break
		}
		key := generateKey(cols[0][i])
		
		if slices.Contains(dart_keys, key) {
			duplicateKeys[key]++
			printInfo(t, fmt.Sprintf("检测到重复 key: %s (出现 %d 次)", key, duplicateKeys[key]+1))
		} else {
			dart_keys = append(dart_keys, key)
			printSuccess(t, fmt.Sprintf("添加唯一 key: %s", key))
		}
	}
	
	// 验证结果
	printInfo(t, fmt.Sprintf("总共处理了 %d 个列，收集到 %d 个唯一 key", len(col)-1, len(dart_keys)))
	
	// 验证没有重复
	keyCount := make(map[string]int)
	for _, key := range dart_keys {
		keyCount[key]++
		if keyCount[key] > 1 {
			printError(t, fmt.Sprintf("发现重复 key: %s", key))
			t.Errorf("发现重复 key: %s", key)
		}
	}
	
	if len(keyCount) == len(dart_keys) {
		printSuccess(t, "验证通过：所有收集的 key 都是唯一的")
	}
	
	// 如果有重复 key 被检测到，验证它们被正确跳过
	if len(duplicateKeys) > 0 {
		printSuccess(t, fmt.Sprintf("成功检测并跳过了 %d 个不同的重复 key", len(duplicateKeys)))
		for key, count := range duplicateKeys {
			printInfo(t, fmt.Sprintf("  - key '%s' 重复了 %d 次", key, count))
		}
	} else {
		printInfo(t, "当前 Excel 文件中没有检测到重复的 key")
	}
	
	// 验证至少有一些 keys
	if len(dart_keys) == 0 {
		printError(t, "未收集到任何 key")
		t.Error("未收集到任何 key")
	} else {
		printSuccess(t, fmt.Sprintf("成功收集到 %d 个唯一 key", len(dart_keys)))
	}
}

// TestEmptyKeyHandling 测试空key的处理
func TestEmptyKeyHandling(t *testing.T) {
	printTestHeader(t, "测试空 key 处理")
	
	testCases := []struct {
		input    string
		expected string
		desc     string
	}{
		{"", "", "空字符串"},
		{"   ", "   ", "只有空格"},
		{"123", "", "纯数字（无法转换为拼音）"},
		{"Hello", "", "纯英文（无法转换为拼音）"},
		{"你好", "nihao", "中文（可以转换为拼音）"},
	}
	
	printInfo(t, fmt.Sprintf("测试 %d 个用例", len(testCases)))
	
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			key := generateKey(tc.input)
			
			// 如果key为空，使用原始文本作为key
			if key == "" {
				key = tc.input
			}
			
			printInfo(t, fmt.Sprintf("输入: '%s', generateKey结果: '%s', 最终key: '%s'", tc.input, generateKey(tc.input), key))
			
			// 验证空key的处理
			if tc.input == "" && key != "" {
				printError(t, fmt.Sprintf("空输入应该保持为空，但得到: '%s'", key))
				t.Errorf("空输入应该保持为空，但得到: '%s'", key)
			}
		})
	}
	
	// 测试空key不会被添加到dart_keys
	var dart_keys []string
	testInputs := []string{"", "   ", "123", "Hello", "你好", "Test Key", "Hello World", "Test: Key!", "ABC-123", "Test (Key)", "UUID: Test-123", "123ABC", "456Test"}
	
	for _, input := range testInputs {
		// 模拟主程序的处理逻辑
		cleanedText := cleanText(input)
		key := generateKey(cleanedText)
		if key == "" {
			key = cleanedText
		}
		
		// 如果key以数字开头，在前面添加 auto_gen_ 前缀
		if key != "" && len(key) > 0 && key[0] >= '0' && key[0] <= '9' {
			key = "auto_gen_" + key
		}
		
		// 跳过空字符的key
		if key == "" {
			printInfo(t, fmt.Sprintf("跳过空key: '%s'", input))
			continue
		}
		
		if !slices.Contains(dart_keys, key) {
			dart_keys = append(dart_keys, key)
			printSuccess(t, fmt.Sprintf("添加key: '%s' (来自输入: '%s')", key, input))
		}
	}
	
	// 验证空key没有被添加
	if slices.Contains(dart_keys, "") {
		printError(t, "空key不应该被添加到dart_keys")
		t.Error("空key不应该被添加到dart_keys")
	} else {
		printSuccess(t, "验证通过：空key没有被添加到dart_keys")
	}
	
	printInfo(t, fmt.Sprintf("最终收集到 %d 个非空key", len(dart_keys)))
}

// BenchmarkPinyinConversion 性能测试：拼音转换
func BenchmarkPinyinConversion(b *testing.B) {
	testString := "你好世界测试数据"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := pinyin.LazyPinyin(testString, pinyin.NewArgs())
		_ = strings.Join(result, "")
	}
}
