package biz

import (
	"Lottery/modules/model"
	"context"
	"errors"
	"regexp"
	"time"
)

type CreateUserStorage interface {
	FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.Player, error)
	CreateUser(ctx context.Context, data *model.PlayerCreation) error
}

type createUserBusiness struct {
	store CreateUserStorage
}

func NewCreateUserBusiness(store CreateUserStorage) *createUserBusiness {
	return &createUserBusiness{store: store}
}

func validateFullName(fullName string) error {
	if fullName == "" {
		return errors.New("tên không được để trống!")
	}

	re := regexp.MustCompile(`^[a-zA-Z ]+$`)
	/*re := regexp.MustCompile(`^[a-zA-ZăâđêôơưÁÀẢÃẠĂẰẮẲẴẶÂẦẤẨẪẬĐÈÉẺẼẸÊỀẾỂỄỆÌÍỈĨỊÒÓỎÕỌÔỒỐỔỖỘƠỜỚỞỠỢÙÚỦŨỤƯỪỨỬỮỰÝỲỶỸỴ ]+$`)*/

	if !re.MatchString(fullName) {
		return errors.New("Tên không được có dấu hoặc chứa kí tự đặc biệt và số")
	}
	return nil
}

func validateDate(dateString string) error {
	if dateString == "" {
		return errors.New("ngày sinh không được để trống!")
	}

	re := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

	if !re.MatchString(dateString) {
		return errors.New("ngày sinh phải theo thứ tự năm-tháng-ngày")
	}

	parsedDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return errors.New("ngày sinh không hợp lệ")
	}
	if parsedDate.Year() > 2007 {
		return errors.New("bạn chưa đủ điều kiện pháp luật để đặt mua")
	}

	return nil
}
func validatePhone(phone string) error {
	if phone == "" {
		return errors.New("số điện thoại không được để trống!")
	}
	re := regexp.MustCompile("^0\\d{9,10}$")
	if !re.MatchString(phone) {
		return errors.New("số điện thoại không đúng định dạng - phải bắt đầu từ 0 và có từ 10 dến 11 số")
	}
	return nil
}
func (business *createUserBusiness) CreateNewUser(ctx context.Context, data *model.PlayerCreation) error {
	errPhone := validatePhone(data.PhoneNumber)
	if errPhone != nil {
		return errPhone
	}
	errName := validateFullName(data.NamePlayer)
	if errName != nil {
		return errName
	}
	errDate := validateDate(data.BirthdayDate)
	if errDate != nil {
		return errDate
	}

	user, _ := business.store.FindUser(ctx, map[string]interface{}{"phone_number": data.PhoneNumber})
	if user != nil {
		return model.ErrRegisteredPhoneNumber
	}
	if err := business.store.CreateUser(ctx, data); err != nil {
		return err
	}
	return nil
}
