package biz

import (
	"Lottery/common"
	"Lottery/modules/model"
	"context"
	"errors"
	"regexp"
)

type LoginStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.Player, error)
}

type loginBusiness struct {
	storeUser LoginStorage
}

func NewLoginBusiness(storeUser LoginStorage) *loginBusiness {
	return &loginBusiness{storeUser: storeUser}
}
func validatePhoneLogin(phonenumber string) error {
	if phonenumber == "" {
		return errors.New("số điện thoại không được để trống!")
	}
	re := regexp.MustCompile("^0\\d{9,10}$")
	if !re.MatchString(phonenumber) {
		return errors.New("số điện thoại không đúng định dạng - phải bắt đầu từ 0 và có từ 10 dến 11 số")
	}
	return nil
}

/*
	func (business *loginBusiness) Login(ctx context.Context, data *model.PlayerLogin) error {
		errPhoneLogin := validatePhoneLogin(data.PhoneNumber)
		if errPhoneLogin != nil {
			return errPhoneLogin
		}
		_, err := business.storeUser.FindUser(ctx, map[string]interface{}{"phone_number": data.PhoneNumber})
		if err != nil {
			return model.ErrLoginPhoneNumber
		}
		return nil
	}
*/
func (business *loginBusiness) Login(ctx context.Context, data *model.PlayerLogin) (*model.Player, error) {
	errPhoneLogin := validatePhoneLogin(data.PhoneNumber)
	if errPhoneLogin != nil {
		return nil, errPhoneLogin
	}

	player, err := business.storeUser.FindUser(ctx, map[string]interface{}{"phone_number": data.PhoneNumber})
	if err != nil {
		if errors.Is(err, common.ErrRecordNotFound) {
			// Thông báo chi tiết về lỗi, ví dụ: Số điện thoại không tồn tại
			return nil, model.ErrLoginPhoneNumber
		}
		// Xử lý các lỗi khác nếu cần thiết

		return nil, err
	}
	return player, nil
}
