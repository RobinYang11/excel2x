# Excel 转 JSON 和 Dart 自动化脚本

**GitHub 仓库**: [https://github.com/RobinYang11/excel2x.git](https://github.com/RobinYang11/excel2x.git)

## 项目描述

这是一个Go语言编写的自动化工具，用于将Excel文件转换为JSON文件和Dart配置文件。特别适合需要国际化配置和多语言支持的项目。

**主要功能：**
- 📊 读取Excel表格数据
- 🔤 使用拼音库自动转换中文key为拼音
- 📄 为每一行数据生成独立的JSON文件
- 🎯 自动生成Dart语言的LocaleKeys类文件
- 🔧 跨平台编译支持（macOS、Windows、Linux）
- 📦 内置版本信息和构建时间

---

## 快速开始

### 下载可执行文件

从 `dist/` 目录选择对应平台的可执行文件：

- **macOS Intel**: `excel2json_macos_intel`
- **macOS Apple Silicon (M1/M2/M3)**: `excel2json_macos_arm64`
- **Windows**: `excel2json.exe`
- **Linux**: `excel2json_linux`

### 基本使用

```bash
# macOS/Linux
./excel2json_macos_arm64 -file=./assets/abc.xlsx

# Windows
excel2json.exe -file=./assets/abc.xlsx
```

---

## 参数说明

### 必需参数

无（所有参数都有默认值）

### 可选参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `-file` | - | `./assets/abc.xlsx` | 输入的Excel文件路径 |
| `-output_path` | - | `./` | JSON文件的输出目录路径 |
| `-dart_file_path` | - | `./` | Dart文件输出目录（文件名固定为 `locale_keys.dart`） |
| `-version` | - | `false` | 显示版本信息、构建时间和Git提交哈希 |

### 参数详细说明

#### `-file`
指定要转换的Excel文件路径。文件必须是 `.xlsx` 格式，且包含名为 `Sheet1` 的工作表。

**示例：**
```bash
./excel2json_macos_arm64 -file=/path/to/your/file.xlsx
./excel2json_macos_arm64 -file=./data/translations.xlsx
```

#### `-output_path`
指定生成的JSON文件的输出目录。如果目录不存在，程序会尝试创建（需要权限）。

**示例：**
```bash
./excel2json_macos_arm64 -output_path=./output/
./excel2json_macos_arm64 -output_path=/tmp/json_files/
```

#### `-dart_file_path`
指定生成的Dart LocaleKeys类文件的输出目录。文件名固定为 `locale_keys.dart`。

**示例：**
```bash
./excel2json_macos_arm64 -dart_file_path=./lib
./excel2json_macos_arm64 -dart_file_path=./generated
```

#### `-version`
显示程序的版本信息，包括版本号、构建时间和Git提交哈希。

**示例：**
```bash
./excel2json_macos_arm64 -version
```

**输出示例：**
```
Version: 1.0.0
Build Time: 2026-01-16 10:30:45
Git Commit: a1b2c3d
```

---

## 使用示例

### 示例 1：使用默认参数

```bash
# 使用默认参数（assets 目录下的 abc.xlsx）
./excel2json_macos_arm64
```

### 示例 2：指定输入文件和输出目录

```bash
# 指定Excel文件和输出目录
./excel2json_macos_arm64 -file=./data/translations.xlsx -output_path=./output/
```

### 示例 3：完整参数配置

```bash
# 指定所有参数
./excel2json_macos_arm64 \
  -file=./data/translations.xlsx \
  -output_path=./output/json/ \
  -dart_file_path=./lib
```

### 示例 4：Windows 使用

```bash
# Windows 命令行
excel2json.exe -file=C:\data\translations.xlsx -output_path=C:\output\
```

---

## Excel 文件格式要求

### 文件结构

Excel文件必须包含名为 `Sheet1` 的工作表，格式要求如下：

1. **第一行**：作为输出JSON文件的文件名（不含扩展名），**必须是英文标识**（如：zh、en、ja等）
2. **第一列**：作为JSON的key，支持中文，程序会自动转换为拼音
3. **数据区域**：从第二行第二列开始，每行对应一个JSON文件

### 示例 Excel 结构

| （第一列，用于生成key） | zh（第一行，JSON文件名） | en | ja |
|---------------------|----------------------|----|----|
| 中文标题 | 你好 | Hello | こんにちは |
| 世界 | 世界 | World | 世界 |
| 测试 | 测试 | Test | テスト |

**说明：**
- **第一行**：必须是英文标识（如 `zh`、`en`、`ja`），用于生成 JSON 文件名（如 `zh.json`、`en.json`）
- **第一列**：是中文标题，用于生成 JSON 的 key（程序会自动转换为拼音，去掉空格和特殊字符）

### 转换结果

**生成的JSON文件：**

`zh.json`:
```json
{
  "zhongwenbiaoti": "你好",
  "shijie": "世界",
  "ceshi": "测试"
}
```

`en.json`:
```json
{
  "zhongwenbiaoti": "Hello",
  "shijie": "World",
  "ceshi": "Test"
}
```

`ja.json`:
```json
{
  "zhongwenbiaoti": "こんにちは",
  "shijie": "世界",
  "ceshi": "テスト"
}
```

**生成的Dart文件 (`locale_keys.dart`):**
```dart
class LocaleKeys {
  static const zhongwenbiaoti = 'zhongwenbiaoti';
  static const shijie = 'shijie';
  static const ceshi = 'ceshi';
}
```

---

## 输出说明

### JSON 文件

- 文件名：使用Excel第一列的值作为文件名（`.json` 扩展名）
- 内容：每行数据转换为一个JSON对象，key为第一行表头（中文转拼音），value为对应单元格的值
- 位置：保存在 `-output_path` 指定的目录中

### Dart 文件

- 文件名：`locale_keys.dart`（固定文件名，不可自定义）
- 输出目录：通过 `-dart_file_path` 参数指定
- 内容：包含所有表头转换后的拼音key，作为静态常量
- 用途：用于Flutter/Dart项目的国际化配置

---

## 常见问题

### Q: 程序提示找不到Excel文件？
A: 检查 `-file` 参数指定的路径是否正确，使用绝对路径或相对路径都可以。

### Q: 生成的JSON文件key是拼音，如何自定义？
A: 目前程序自动将中文表头转换为拼音。如需自定义key，建议在Excel第一行直接使用英文或拼音。

### Q: 支持哪些Excel格式？
A: 目前仅支持 `.xlsx` 格式（Excel 2007及以上版本）。

### Q: 如何查看程序版本？
A: 使用 `-version` 参数：`./excel2json_macos_arm64 -version`

---

## 开发说明

### 环境要求

- Go 1.16 或更高版本
- 依赖包：
  - `github.com/mozillazg/go-pinyin` - 中文转拼音
  - `github.com/xuri/excelize/v2` - Excel文件处理

### 安装依赖

```bash
go mod download
```

### 本地开发运行

```bash
go run main.go -file=./assets/abc.xlsx -output_path=./output/
```

### 构建可执行文件

使用提供的构建脚本：

```bash
chmod +x build.sh
./build.sh
```

构建完成后，可执行文件将保存在 `./dist/` 目录中。

### 手动构建

```bash
# macOS
GOOS=darwin GOARCH=arm64 go build -o excel2json main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o excel2json.exe main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o excel2json main.go
```

