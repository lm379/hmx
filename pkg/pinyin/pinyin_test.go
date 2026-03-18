package pinyin

import (
	"testing"
)

func TestTokenizeChinese(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
		{
			name:     "Simple Chinese characters",
			text:     "黄梅戏",
			expected: "黄 梅 戏",
		},
		{
			name:     "Chinese with numbers",
			text:     "黄梅戏123",
			expected: "黄 梅 戏123",
		},
		{
			name:     "Chinese with English",
			text:     "黄梅戏HMX",
			expected: "黄 梅 戏HMX",
		},
		{
			name:     "Chinese with symbols",
			text:     "黄梅戏-戏曲",
			expected: "黄 梅 戏- 戏 曲",
		},
		{
			name:     "Single Chinese character",
			text:     "黄",
			expected: "黄",
		},
		{
			name:     "Two Chinese characters",
			text:     "黄梅",
			expected: "黄 梅",
		},
		{
			name:     "Mixed text with Chinese and numbers",
			text:     "这是一个测试2024年",
			expected: "这 是 一 个 测 试2024 年",
		},
		{
			name:     "Chinese with punctuation",
			text:     "黄梅戏，戏曲。",
			expected: "黄 梅 戏， 戏 曲。",
		},
		{
			name:     "Only English",
			text:     "Hello World",
			expected: "Hello World",
		},
		{
			name:     "Only numbers",
			text:     "123456",
			expected: "123456",
		},
		{
			name:     "Only symbols",
			text:     "!@#$%",
			expected: "!@#$%",
		},
		{
			name:     "Chinese with spaces",
			text:     "黄梅 戏",
			expected: "黄 梅  戏",
		},
		{
			name:     "Mixed content",
			text:     "用户user123",
			expected: "用 户user123",
		},
		{
			name:     "Long Chinese text",
			text:     "中华人民共和国",
			expected: "中 华 人 民 共 和 国",
		},
		{
			name:     "Chinese with emoji",
			text:     "黄梅戏😊",
			expected: "黄 梅 戏😊",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TokenizeChinese(tt.text)
			if result != tt.expected {
				t.Errorf("TokenizeChinese() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestToPinyin(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: "",
		},
		{
			name:     "Simple Chinese",
			text:     "黄梅戏",
			expected: "huangmeixi",
		},
		{
			name:     "Single character",
			text:     "黄",
			expected: "huang",
		},
		{
			name:     "Chinese with English",
			text:     "黄梅戏HMX",
			expected: "huangmeixihmx",
		},
		{
			name:     "Chinese with numbers",
			text:     "黄梅戏123",
			expected: "huangmeixi123",
		},
		{
			name:     "Chinese with symbols",
			text:     "黄梅戏-戏曲",
			expected: "huangmeixi-xiqu",
		},
		{
			name:     "Only English",
			text:     "Hello",
			expected: "hello",
		},
		{
			name:     "Mixed case English",
			text:     "HelloWorld",
			expected: "helloworld",
		},
		{
			name:     "Numbers only",
			text:     "12345",
			expected: "12345",
		},
		{
			name:     "Mixed Chinese and English",
			text:     "用户user",
			expected: "yonghuuser",
		},
		{
			name:     "Chinese with punctuation",
			text:     "黄梅戏，戏曲",
			expected: "huangmeixi，xiqu",
		},
		{
			name:     "Special characters",
			text:     "!@#$%",
			expected: "!@#$%",
		},
		{
			name:     "Mixed content",
			text:     "测试Test123",
			expected: "ceshitest123",
		},
		{
			name:     "Mixed case Chinese conversion",
			text:     "黄梅戏",
			expected: "huangmeixi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPinyin(tt.text)
			if result != tt.expected {
				t.Errorf("ToPinyin() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestToPinyinSlice(t *testing.T) {
	tests := []struct {
		name     string
		texts    []string
		expected []string
	}{
		{
			name:     "Empty slice",
			texts:    []string{},
			expected: []string{},
		},
		{
			name:     "Single element",
			texts:    []string{"黄梅戏"},
			expected: []string{"huangmeixi"},
		},
		{
			name:     "Multiple elements",
			texts:    []string{"黄梅戏", "戏曲", "艺术"},
			expected: []string{"huangmeixi", "xiqu", "yishu"},
		},
		{
			name:     "Mixed content",
			texts:    []string{"黄", "戏", "123", "Hello"},
			expected: []string{"huang", "xi", "123", "hello"},
		},
		{
			name:     "With empty strings",
			texts:    []string{"黄梅戏", "", "戏曲"},
			expected: []string{"huangmeixi", "", "xiqu"},
		},
		{
			name:     "Special characters",
			texts:    []string{"!@#", "测试"},
			expected: []string{"!@#", "ceshi"},
		},
		{
			name:     "Long slice",
			texts:    []string{"黄", "梅", "戏", "艺", "术"},
			expected: []string{"huang", "mei", "xi", "yi", "shu"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToPinyinSlice(tt.texts)
			if len(result) != len(tt.expected) {
				t.Errorf("ToPinyinSlice() length = %d, want %d", len(result), len(tt.expected))
				return
			}

			for i := range tt.expected {
				if result[i] != tt.expected[i] {
					t.Errorf("ToPinyinSlice()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestTokenizeChineseEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Repeating characters",
			text:     "哈哈哈哈",
			expected: "哈 哈 哈 哈",
		},
		{
			name:     "Chinese with tabs",
			text:     "黄\t梅\t戏",
			expected: "黄\t 梅\t 戏",
		},
		{
			name:     "Chinese with newlines",
			text:     "黄\n梅\n戏",
			expected: "黄\n 梅\n 戏",
		},
		{
			name:     "Chinese with multiple consecutive Chinese",
			text:     "黄梅戏曲艺术",
			expected: "黄 梅 戏 曲 艺 术",
		},
		{
			name:     "Chinese mixed with multiple consecutive English",
			text:     "黄HelloWorld梅",
			expected: "黄HelloWorld 梅",
		},
		{
			name:     "All special characters",
			text:     "！@#￥%……&*（）",
			expected: "！@#￥%……&*（）",
		},
		{
			name:     "Chinese with mixed numbers",
			text:     "2024年5月1日",
			expected: "2024 年5 月1 日",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TokenizeChinese(tt.text)
			if result != tt.expected {
				t.Errorf("TokenizeChinese() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestToPinyinConsistency(t *testing.T) {
	// Test that same input produces consistent output
	input := "黄梅戏"

	result1 := ToPinyin(input)
	result2 := ToPinyin(input)

	if result1 != result2 {
		t.Errorf("ToPinyin() should produce consistent results, got %q and %q", result1, result2)
	}
}

func TestPinyinFunctionsNilSafe(t *testing.T) {
	// Test empty string handling
	tokenizeResult := TokenizeChinese("")
	if tokenizeResult != "" {
		t.Errorf("TokenizeChinese(\"\") should return empty string, got %q", tokenizeResult)
	}

	pinyinResult := ToPinyin("")
	if pinyinResult != "" {
		t.Errorf("ToPinyin(\"\") should return empty string, got %q", pinyinResult)
	}
}

func TestToPinyinSliceLength(t *testing.T) {
	inputs := []string{"黄", "梅", "戏", "艺", "术"}
	results := ToPinyinSlice(inputs)

	if len(results) != len(inputs) {
		t.Errorf("ToPinyinSlice() returned %d results, expected %d", len(results), len(inputs))
	}
}

func TestMixedContent(t *testing.T) {
	tests := []struct {
		name             string
		text             string
		tokenizeExpected string
		pinyinExpected   string
	}{
		{
			name:             "Chinese + English + Numbers",
			text:             "黄梅戏HMX2024",
			tokenizeExpected: "黄 梅 戏HMX2024",
			pinyinExpected:   "huangmeixihmx2024",
		},
		{
			name:             "Email format",
			text:             "用户@example.com",
			tokenizeExpected: "用 户@example.com",
			pinyinExpected:   "yonghu@example.com",
		},
		{
			name:             "URL format",
			text:             "黄梅戏.com",
			tokenizeExpected: "黄 梅 戏.com",
			pinyinExpected:   "huangmeixi.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizeResult := TokenizeChinese(tt.text)
			if tokenizeResult != tt.tokenizeExpected {
				t.Errorf("TokenizeChinese() = %q, want %q", tokenizeResult, tt.tokenizeExpected)
			}

			pinyinResult := ToPinyin(tt.text)
			if pinyinResult != tt.pinyinExpected {
				t.Errorf("ToPinyin() = %q, want %q", pinyinResult, tt.pinyinExpected)
			}
		})
	}
}
