package service

import (
	"fmt"
	"os"

	"chrome/storage"
	"chrome/utils"
)

type ChromeService struct {
	outputFile string
}

func NewChromeService(outputFile string) *ChromeService {
	return &ChromeService{
		outputFile: outputFile,
	}
}

func (s *ChromeService) ExportPasswords() error {
	// 获取Chrome路径
	loginDataPath, localStatePath := utils.GetChromePaths()

	// 复制数据库文件
	tempDB := "Login Data.db"
	if err := utils.CopyFile(loginDataPath, tempDB); err != nil {
		return fmt.Errorf("复制数据库失败: %v", err)
	}
	defer os.Remove(tempDB)

	// 获取主密钥
	masterKey, err := utils.GetMasterKey(localStatePath)
	if err != nil {
		return fmt.Errorf("获取主密钥失败: %v", err)
	}

	// 打开数据库
	db, err := storage.NewChromeDB(tempDB)
	if err != nil {
		return fmt.Errorf("打开数据库失败: %v", err)
	}
	defer db.Close()

	// 获取登录数据
	results, err := db.GetLoginData(masterKey)
	if err != nil {
		return fmt.Errorf("获取登录数据失败: %v", err)
	}

	// 保存结果
	if err := utils.SaveResults(s.outputFile, results); err != nil {
		return fmt.Errorf("保存结果失败: %v", err)
	}

	return nil
}
