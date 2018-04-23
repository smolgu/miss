package vk

import (
	"fmt"
	"net/url"

	"github.com/smolgu/miss/pkg/errors"

	"github.com/zhuharev/vkutil"
)

// GetUser returns info about user by id
func GetUser(token string, vkID int64) (user vkutil.User, err error) {
	var (
		u   = vkutil.New()
		res []vkutil.User
	)
	u.SetDebug(true)
	const fields = "photo_200,sex,followers_count,counters"
	u.VkApi.AccessToken = token
	u.VkApi.Lang = "ru"
	params := url.Values{}
	params.Set("fields", fields)
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
	return GetUser(token, 0)
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
