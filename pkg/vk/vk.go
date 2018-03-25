package vk

import (
	"fmt"
	"net/url"

	"github.com/zhuharev/vkutil"
)

func userGet(token string, withAvatarURL bool) (user vkutil.User, err error) {
	var (
		u   = vkutil.New()
		res []vkutil.User
	)
	u.VkApi.AccessToken = token
	u.VkApi.Lang = "ru"
	params := url.Values{}
	if withAvatarURL {
		params.Set("fields", "photo_200")
	}
	res, err = u.UsersGet(nil, params)
	if err != nil {
		return
	}
	if len(res) != 1 {
		err = fmt.Errorf("Token invalid")
		return
	}
	return res[0], nil
}

// CheckToken return user id of vk.com user
func CheckToken(token string) (int64, error) {
	user, err := userGet(token, true)
	if err != nil {
		return 0, err
	}
	return int64(user.Id), nil
}

// GetAvatarURL return user's avatar
func GetAvatarURL(token string) (string, error) {
	u, err := userGet(token, true)
	if err != nil {
		return "", err
	}
	return u.Photo200, nil
}
