// Excel to JSON converter
// Author: RobinYang
// Date: 2026.1.16
// Description: This program reads an Excel file and converts each row into a separate JSON file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/mozillazg/go-pinyin"
	"github.com/xuri/excelize/v2"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

const dart_file_template = `class LocaleKeys {
%s
}
`

// truncateTextBeforePinyin 在转换为拼音之前，如果文本转换为拼音后长度会超过8，则截取前4个和后4个字符
func truncateTextBeforePinyin(text string) string {
	// 先转换为拼音检查长度
	keySlice := pinyin.LazyPinyin(text, pinyin.NewArgs())
	pinyinKey := strings.Join(keySlice, "")

	// 如果拼音长度超过8，则对原始文本截取前4个和后4个字符
	if len(pinyinKey) > 8 {
		runes := []rune(text)
		if len(runes) > 8 {
			first4 := string(runes[:4])
			last4 := string(runes[len(runes)-4:])
			return first4 + "_" + last4
		}
	}
	return text
}

// cleanText 清理文本：去掉空格和特殊字符，只保留字母、数字和中文字符
func cleanText(text string) string {
	// 使用正则表达式，只保留字母、数字和中文字符
	// \p{L} 匹配所有字母（包括中文），\p{N} 匹配所有数字
	// 或者使用 [a-zA-Z0-9] 匹配英文和数字，[\u4e00-\u9fff] 匹配中文
	re := regexp.MustCompile(`[^a-zA-Z0-9\p{Han}]`)
	cleaned := re.ReplaceAllString(text, "")
	return cleaned
}

// generateKey 生成key：先截取文本（如果需要），然后转换为拼音
func generateKey(text string) string {
	truncatedText := truncateTextBeforePinyin(text)
	keySlice := pinyin.LazyPinyin(truncatedText, pinyin.NewArgs())
	key := strings.Join(keySlice, "")
	return key
}

func main() {

	input_file := flag.String("file", "./assets/abc.xlsx", "input excel file")
	output_path := flag.String("output_path", "./", "output json file path")
	dart_file_path := flag.String("dart_file_path", "./", "dart file output directory")
	showVersion := flag.Bool("version", false, "show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Version: %s\nBuild Time: %s\nGit Commit: %s\n", version, buildTime, gitCommit)
		return
	}

	fmt.Println("input_file:", *input_file)
	fmt.Println("output_path:", *output_path)

	f, err := excelize.OpenFile(*input_file)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, col := range cols[0] {
		res := generateKey(col)
		fmt.Println(res)
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	len_of_col := len(cols)
	var dart_keys []string
	// forech first row , create map use first row as key
	for row_index, row := range rows[0] {

		file_data := make(map[string]string, len_of_col-1)

		col := cols[row_index]

		for i := 1; i < len(col); i++ {
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

			// Collect keys for dart file
			if row_index == 0 {
				// 3. 跳过空字符的key
				if key == "" {
					fmt.Printf("跳过空key: %v\n", originalText)
					continue
				}
				// 去重
				if slices.Contains(dart_keys, key) {
					fmt.Println("key already exists:", key)
					continue
				}
				dart_keys = append(dart_keys, key)
			}
			// fmt.Printf("key: %v\n", key)
			file_data[key] = col[i]
		}
		WriteJson(file_data, fmt.Sprintf("%s/%s.json", *output_path, string(row)))
	}

	// Generate dart file
	var dart_content strings.Builder
	for _, key := range dart_keys {
		dart_content.WriteString(fmt.Sprintf("  static const %s = '%s';\n", key, key))
	}
	final_dart := fmt.Sprintf(dart_file_template, strings.TrimRight(dart_content.String(), "\n"))

	dart_file_full_path := filepath.Join(*dart_file_path, "locale_keys.dart")
	dart_file, err := os.Create(dart_file_full_path)
	if err != nil {
		fmt.Println("create dart file error:", err)
		return
	}
	defer dart_file.Close()
	dart_file.WriteString(final_dart)

}

func WriteJson(data map[string]string, name string) {
	fmt.Println("write json file:", name)
	file, err := os.Create(name)

	if err != nil {
		fmt.Println("create file error:", err)
		return
	}
	defer file.Close()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}
	file.WriteString(string(jsonData))
}
