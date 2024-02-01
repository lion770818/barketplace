package entity

import (
	"time"
)

type Transaction struct {
	ID            uint64    `gorm:"primary_key;auto_increment;comment:'流水號 主鍵'" json:"id"`
	UserID        uint64    `gorm:"size:100;not null; comment:'交易用戶' " json:"user_id"`
	TransactionID string    `gorm:"size:100;not null;unique; comment:'交易ID'" json:"transactionId"`                // 交易ID
	Satus         int       `gorm:"size:255;null; comment:'交易狀態 0=idle 1=start 2=fail 3=success';" json:"status"` // 交易狀態
	Description   string    `gorm:"text;not null;comment:'描述'" json:"description"`                                // 描述
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'交易日期'" json:"created_at"`                   // 交易日期
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;comment:'更新日期'" json:"updated_at"`                   // 更新狀態日期
}

func (t *Transaction) Validate() map[string]string {
	var errorMessages = make(map[string]string)

	if t.UserID <= 0 {
		errorMessages["user_id_required"] = "title is required"
	} else if len(t.TransactionID) == 0 {
		errorMessages["transactionId_required"] = "transactionId is required"
	}

	return errorMessages
}
