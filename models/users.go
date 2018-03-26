package models

import (
	"fmt"

	"github.com/golang/protobuf/ptypes"
	pkgErrors "github.com/pkg/errors"
	"github.com/zhuharev/vkutil"

	"github.com/smolgu/miss/pkg/errors"
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
	user := new(User)
	has, err := db.Id(userID).Get(user)
	if err != nil {
		return nil, errors.New(errors.ErrPublicUnknownError, pkgErrors.Wrapf(err, "get user from db user=%d", userID))
	}
	if !has {
		return nil, errors.New(fmt.Errorf("пользователь не найден"), errors.ErrNotFound, errors.ErrNotFound)
	}
	return user, nil
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
		CreatedAt: ptypes.TimestampNow(),
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
ORDER BY random()
LIMIT 10
`, voterID).Find(&res)
	return
}
