package models

import (
	"fmt"

	"github.com/smolgu/miss/pkg/errors"

	pkgErrors "github.com/pkg/errors"
	"github.com/zhuharev/vkutil"
)

type userModel int

var Users userModel

func (u User) TableName() string {
	return "users"
}

func (u User) ObjectType() ObjectType {
	return ObjectType_ObjectUser
}

func (userModel) Get(userID int64) (*User, error) {
	u := new(User)
	has, err := db.Id(userID).Get(&u)
	if err != nil {
		return nil, errors.New(fmt.Errorf("Неизвестная ошибка"), pkgErrors.Wrap(err, "get user from db"))
	}
	if !has {
		return nil, errors.New(fmt.Errorf("Пользователь не найден"), pkgErrors.Wrapf(err, "user not found user_id=%d", userID), errors.ErrNotFound)
	}
	return u, nil
}

func (userModel) GetByVkID(vkID int64) (*User, error) {
	var u User
	has, err := db.Where("vk_id = ?", vkID).Get(&u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.ErrNotFound
	}
	return &u, nil
}

func (userModel) CreateByVKUser(vkUser vkutil.User) (*User, error) {
	user := User{
		FirstName: vkUser.FirstName,
		LastName:  vkUser.LastName,
		VkId:      int64(vkUser.Id),
		AvatarUrl: vkUser.Photo200,
	}
	_, err := db.Insert(&user)
	return &user, err
}

func (userModel) Random(voterID int64) (res []*User, err error) {
	err = db.SQL(`SELECT *
FROM users
WHERE
  id NOT IN (
		SELECT target_id
		FROM votes
		WHERE voter_id = ?
	)
OREDER BY random()
LIMIT 10
`, voterID).Find(&res)
	return
}
