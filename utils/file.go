package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"chrome/storage"

	"github.com/tidwall/gjson"

	"chrome/crypto"
)

// GetChromePaths returns Chrome's login data and local state file paths
func GetChromePaths() (loginDataPath, localStatePath string) {
	userProfile := os.Getenv("USERPROFILE")
	loginDataPath = filepath.Join(userProfile, "AppData", "Local", "Google", "Chrome", "User Data", "Default", "Login Data")
	localStatePath = filepath.Join(userProfile, "AppData", "Local", "Google", "Chrome", "User Data", "Local State")
	return
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// GetMasterKey retrieves and decrypts the master key from Chrome's local state
func GetMasterKey(localStatePath string) ([]byte, error) {
	data, err := os.ReadFile(localStatePath)
	if err != nil {
		return nil, err
	}

	encryptedKey := gjson.Get(string(data), "os_crypt.encrypted_key")
	if !encryptedKey.Exists() {
		return nil, fmt.Errorf("找不到加密密钥")
	}

	key, err := base64.StdEncoding.DecodeString(encryptedKey.String())
	if err != nil {
		return nil, err
	}

	// 修改这里，使用 crypto 包的函数
	return crypto.DecryptWithDPAPI(key[5:])
}

func SaveResults(filename string, results []storage.LoginData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintln(file, "Chrome密码导出结果:")
	fmt.Fprintln(file, "----------------------------------------")

	for _, data := range results {
		fmt.Fprintf(file, "URL: %s\n", data.LoginURL)
		fmt.Fprintf(file, "用户名: %s\n", data.UserName)
		fmt.Fprintf(file, "密码: %s\n", data.Password)
		fmt.Fprintf(file, "创建时间: %s\n", data.CreateDate.Format("2006-01-02 15:04:05"))
		fmt.Fprintln(file, "----------------------------------------")
	}

	return nil
}
