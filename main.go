package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang-collections/collections/stack"
	"github.com/vickxxx/my_dbg/models"
)

const logFile = "/tmp/mysqld.trace"

func main() {
	fmt.Println("ss")
	// 打开文件
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	// 确保在函数结束时关闭文件
	defer file.Close()

	// 创建一个scanner来读取文件
	scanner := bufio.NewScanner(file)
	rootLog := &models.Line{
		Title: "root",
	}
	logStack := stack.New()
	logStack.Push(rootLog)
	// var parentLog *models.Line
	// 逐行读取文件
	for scanner.Scan() {
		// 获取当前行的文本
		line := scanner.Text()
		// 打印当前行
		fmt.Println(line)
		logLine := models.ParseLine(line)

		fmt.Println(logLine.Depth)
		parentLog, ok := logStack.Peek().(*models.Line)
		if !ok {
			parentLog = rootLog
			logStack.Push(rootLog)
		}

		// 补充祖先
		if parentLog.Depth == 0 && logLine.Depth > 1 {
			parentLog = addParents(rootLog, logStack, logLine.Depth)
			parentLog.Detail = append(parentLog.Detail, logLine)
			continue
		}

		switch logLine.Type {
		case 1: // push stack
			parentLog.Detail = append(parentLog.Detail, logLine)

			logStack.Push(logLine)

		case 2: // pop stack

			parentLog, ok := logStack.Peek().(*models.Line)
			if ok {
				if parentLog.Title != logLine.Title {
					parentLog.Title += " - " + logLine.Title
				}
			}
			logStack.Pop()

		default:
			parentLog.Detail = append(parentLog.Detail, logLine)
		}

		// break
	}

	fmt.Println(spew.Sdump(rootLog.Detail))

	fmt.Println(len(rootLog.Detail))
	// 检查是否有读取错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	// 创建 JSON 编码器
	// 创建输出 JSON 文件
	outputFile, err := os.Create("output.json")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ") // 设置缩进以使输出更易读
	// 将数据写入 JSON 文件
	if err := encoder.Encode(rootLog.Detail); err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("JSON data has been written to output.json")
}

func addParents(root *models.Line, stack *stack.Stack, depth int64) *models.Line {

	var parent *models.Line
	for i := range depth {
		if i == 0 {
			continue
		}
		// println(i)
		p := models.Line{
			Depth: i,
			// Type:  1,
			Title: "no root",
		}
		// *stack = append(*stack, &p)
		stack.Push(&p)
		if parent == nil {
			parent = &p
			root.Detail = append(root.Detail, parent)
		} else {
			parent.Detail = append(parent.Detail, &p)
			parent = &p
		}
	}
	// color.Greenln(spew.Sdump(root))
	return parent
}
