package acl

import (
	"crypto/sha256"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

const Super = "admin"

type AdminRole struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Authority  Authority `json:"authority,omitempty"`
	CreateTime time.Time `json:"create_time" gorm:"->"` // 只读
	UpdateTime time.Time `json:"update_time" gorm:"->"` // 只读
}

func (*AdminRole) TableName() string {
	return "admin_role"
}

type AdminUser struct {
	ID         int        `json:"id"`
	Username   string     `json:"username"`
	Password   string     `json:"-"`
	RoleID     int        `json:"role_id"`
	Status     int8       `json:"status"`
	CreateTime time.Time  `json:"create_time" gorm:"->"` // 只读
	UpdateTime time.Time  `json:"update_time" gorm:"->"` // 只读
	AdminRole  *AdminRole `json:"admin_role,omitempty" gorm:"foreignKey:ID;references:RoleID"`
}

func (*AdminUser) TableName() string {
	return "admin_user"
}

func (a *AdminUser) BeforeSave(*gorm.DB) error {
	if a.Password != "" {
		a.Password = cryptoHash(a.Password)
	}
	return nil
}

type Authority map[string]int8

func (auth *Authority) Scan(value any) error {
	if value == nil {
		return nil
	}
	b := value.([]byte)
	return json.Unmarshal(b, auth) // receiver必须为指针
}

func (auth Authority) Value() (driver.Value, error) {
	if auth == nil {
		return []byte{'{', '}'}, nil
	}
	return json.Marshal(auth) // receiver不能为指针
}

func CheckPassword(input, crypt string) bool {
	return cryptoHash(input) == crypt
}

func cryptoHash(pwd string) string {
	h := sha256.Sum256([]byte(pwd))
	return base64.URLEncoding.EncodeToString(h[:30])
}

type AdminToken struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Authority Authority `json:"authority"`
}
