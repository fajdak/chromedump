package main

import (
	"fmt"
	"os"

	"chrome/service" // 修改这里的导入路径
)

func main() {
	chromeService := service.NewChromeService("chrome_passwords.txt")
	if err := chromeService.ExportPasswords(); err != nil {
		fmt.Printf("导出密码失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("密码已成功导出到 chrome_passwords.txt")
}
