package output

import (
	"encoding/json"
	"fmt"

	"github.com/fatih/color"
)

type Output interface {
	Success(msg string)
	Error(err error)
	Print(data interface{})
}

// HumanOutput 人类可读输出
type HumanOutput struct{}

func (h *HumanOutput) Success(msg string) {
	color.New(color.FgGreen).Println("✓", msg)
}

func (h *HumanOutput) Error(err error) {
	color.New(color.FgRed).Println("✗", err.Error())
}

func (h *HumanOutput) Print(data interface{}) {
	fmt.Printf("%#v\n", data)
}

// JSONOutput JSON 输出
type JSONOutput struct{}

func (j *JSONOutput) Success(msg string) {
	j.Print(map[string]string{"status": "success", "message": msg})
}

func (j *JSONOutput) Error(err error) {
	j.Print(map[string]string{"status": "error", "error": err.Error()})
}

func (j *JSONOutput) Print(data interface{}) {
	jsonBytes, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(jsonBytes))
}

// NewOutput 创建输出器
func NewOutput(isJSON bool) Output {
	if isJSON {
		return &JSONOutput{}
	}
	return &HumanOutput{}
}
