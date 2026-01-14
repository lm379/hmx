package pinyin

import (
	"strings"
	"unicode"

	"github.com/mozillazg/go-pinyin"
)

var pinyinArgs = pinyin.NewArgs()

func init() {
	// 配置拼音转换参数
	pinyinArgs.Style = pinyin.Normal // 不带声调
	pinyinArgs.Heteronym = false     // 不返回多音字
	pinyinArgs.Separator = ""        // 拼音之间不加分隔符
	pinyinArgs.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)} // 非汉字字符保持原样
	}
}

// TokenizeChinese 将中文字符串分词（用空格分隔每个字符）
// 这样可以让 Meilisearch 更好地处理模糊搜索
func TokenizeChinese(text string) string {
	if text == "" {
		return ""
	}

	var result strings.Builder
	runes := []rune(text)
	
	for i, r := range runes {
		if i > 0 && unicode.Is(unicode.Han, r) {
			result.WriteRune(' ') // 中文字符前加空格
		}
		result.WriteRune(r)
	}

	return result.String()
}

// ToPinyin 将中文转换为拼音
func ToPinyin(text string) string {
	if text == "" {
		return ""
	}

	result := pinyin.Pinyin(text, pinyinArgs)
	var pinyinStr strings.Builder

	for _, py := range result {
		if len(py) > 0 {
			pinyinStr.WriteString(py[0])
		}
	}

	return strings.ToLower(pinyinStr.String())
}

// ToPinyinSlice 将字符串切片转换为拼音切片
func ToPinyinSlice(texts []string) []string {
	result := make([]string, len(texts))
	for i, text := range texts {
		result[i] = ToPinyin(text)
	}
	return result
}
