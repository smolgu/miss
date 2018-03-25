package models

import "ug/errors"

type userModel int

var Users userModel

func (u User) TableName() string {
	return "users"
}

func (u User) ObjectType() ObjectType {
	return ObjectType_ObjectUser
}

func (userModel) Get(userID int64) (*User, error) {
	return nil, nil
}

func (userModel) GetByVkID(vkID int64) (*User, error) {
	var u User
	has, err := db.Where("vk_id = ?", vkID).Get(&u)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.NotFound
	}
	return nil, nil
}

func (um userModel) GetByVkIDOrRegister(vkID int64) (*User, error) {
	user, err := um.GetByVkID(vkID)
	if err != nil {
		if err != errors.NotFound {
			return nil, err
		}
	}
	_ = user

	return nil, nil
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
