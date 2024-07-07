package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tealeg/xlsx"
)

type Problem struct {
	Question string
	Options  []string
	Answer   string
	Type     string // 存储题目类型
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// 判断题目类型的函数
func determineType(answer string) string {
	if len(answer) == 1 {
		if answer == "Y" || answer == "N" {
			return "判断题"
		}
		return "单选题"
	}
	return "多选题"
}

func main() {
	// 创建一个新的随机数生成器
	randSource := rand.NewSource(time.Now().UnixNano())
	randGenerator := rand.New(randSource)

	// 打开 Excel 文件
	xlFile, err := xlsx.OpenFile("problems.xlsx")
	if err != nil {
		log.Fatalf("无法打开 Excel 文件: %s\n", err)
	}

	// 读取第一个工作表
	sheet := xlFile.Sheets[0]

	// 解析题目
	problems := make([]Problem, 0)

	// 跳过第一行，从第二行开始读取
	for i, row := range sheet.Rows {
		if i == 0 {
			continue // 跳过第一行
		}
		question := row.Cells[0].String()
		options := []string{
			row.Cells[1].String(),
			row.Cells[2].String(),
			row.Cells[3].String(),
			row.Cells[4].String(),
		}
		answer := row.Cells[5].String()
		problemType := determineType(strings.TrimSpace(answer)) // 判断题目类型

		problem := Problem{
			Question: question,
			Options:  options,
			Answer:   answer,
			Type:     problemType,
		}

		problems = append(problems, problem)
	}

	totalProblems := len(problems)

	// 初始化欢迎信息
	fmt.Println("题库包含：单选、多选、判断")
	fmt.Println("BY-CCSU_YZT")
	fmt.Println("按任意键继续...")
	fmt.Scanln()
	correct, allnum := 0, 0
	for totalProblems > 0 {
		// 打印未答题目数量
		clearScreen()
		fmt.Printf("当前还有 %d 题未答。当前正确率是 %d / %d \n", totalProblems, correct, allnum)
		allnum++
		// 随机选择一个未答过的题目
		problemIndex := randGenerator.Intn(len(problems))
		selectedProblem := problems[problemIndex]

		// 展示题目和类型
		fmt.Println("题目类型:", selectedProblem.Type)
		fmt.Println("题目:", selectedProblem.Question)

		if selectedProblem.Type == "判断题" {
			fmt.Println("选项:")
			fmt.Println("Y) 正确")
			fmt.Println("N) 错误")
			fmt.Print("请输入您的答案 (Y/N): ")
		} else {
			fmt.Println("选项:")
			for i, option := range selectedProblem.Options {
				if option != "" {
					fmt.Printf("%c) %s\n", 'A'+i, option)
				}
			}
			fmt.Print("请输入您的答案 (A/B/C/D): ")
		}

		// 接受用户答案
		var userAnswer string
		fmt.Scanln(&userAnswer)

		// 检查答案是否正确
		correctAnswer := strings.TrimSpace(selectedProblem.Answer)
		if strings.EqualFold(userAnswer, correctAnswer) {
			fmt.Println("回答正确!")
			// 从 problems 中删除已经回答正确的题目
			problems = append(problems[:problemIndex], problems[problemIndex+1:]...)
			totalProblems--
			correct++
		} else {
			fmt.Printf("回答错误。正确答案是 %s.\n", correctAnswer)
		}

		fmt.Println()
	}

	// 完成所有题目后，等待用户按任意键退出
	fmt.Println("恭喜您，所有题目已答完!")
	fmt.Printf("您的正确率是 %d / %d \n", correct, allnum)
	fmt.Println("按任意键退出...")
	fmt.Scanln()
}
