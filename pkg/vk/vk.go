package vk

import (
	"fmt"
	"net/url"

	"github.com/zhuharev/vkutil"

	"github.com/smolgu/miss/pkg/errors"
)

// GetUser returns info about user by id
func GetUser(token string, vkID int64, withAvatars ...bool) (user vkutil.User, err error) {
	var (
		u             = vkutil.New()
		res           []vkutil.User
		withAvatarURL = len(withAvatars) > 0 && withAvatars[0]
	)
	u.VkApi.AccessToken = token
	u.VkApi.Lang = "ru"
	params := url.Values{}
	if withAvatarURL {
		params.Set("fields", "photo_200")
	}
	res, err = u.UsersGet(vkID, params)
	if err != nil {
		err = errors.New(fmt.Errorf("ошибка соединения с ВКонтакте. Попробуйте позже"),
			err)
		return
	}
	if len(res) != 1 {
		err = errors.New(fmt.Errorf("ошибка авторизации во ВКонтакте. Попробуйте позже"),
			err)
		return
	}
	return res[0], nil
}

// UserGetByToken get info about user by token
func UserGetByToken(token string, withAvatarURL bool) (user vkutil.User, err error) {
	return GetUser(token, 0, withAvatarURL)
}

// CheckToken return user id of vk.com user
func CheckToken(token string) (int64, error) {
	user, err := UserGetByToken(token, true)
	if err != nil {
		return 0, err
	}
	return int64(user.Id), nil
}

// GetAvatarURL return user's avatar
func GetAvatarURL(token string) (string, error) {
	u, err := UserGetByToken(token, true)
	if err != nil {
		return "", err
	}
	return u.Photo200, nil
}
