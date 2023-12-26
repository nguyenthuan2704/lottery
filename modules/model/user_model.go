package model

import (
	"Lottery/common"
	"errors"
)

var (
	ErrRegisteredPhoneNumber = common.NewCustomError(
		errors.New("số điện thoại đã được đăng ký"),
		"số điện thoại đã được đăng ký",
		"ErrRegisteredPhoneNumber",
	)
	ErrLoginPhoneNumber = common.NewCustomError(
		errors.New("số điện thoại chưa có đăng ký tài khoản"),
		"số điện thoại chưa có đăng ký tài khoản",
		"ErrLoginPhoneNumber",
	)
)

type Player struct {
	common.SQLModel
	PhoneNumber  string `json:"phone_number" gorm:"column:phone_number;"`
	NamePlayer   string `json:"name" gorm:"column:name;"`
	BirthdayDate string `json:"date_of_birth" gorm:"column:date_of_birth;"`
}

func (Player) TableName() string { return "players" }

type PlayerCreation struct {
	common.SQLModel
	PhoneNumber  string `json:"phone_number" gorm:"column:phone_number;"`
	NamePlayer   string `json:"name" gorm:"column:name;"`
	BirthdayDate string `json:"date_of_birth" gorm:"column:date_of_birth;"`
}

func (PlayerCreation) TableName() string { return Player{}.TableName() }

type PlayerLogin struct {
	common.SQLModel
	PhoneNumber string `json:"phone_number" gorm:"column:phone_number;"`
}

func (PlayerLogin) TableName() string {
	return Player{}.TableName()
}
