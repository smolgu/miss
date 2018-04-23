// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/smolgu/miss/pkg/setting"
	"github.com/smolgu/miss/pkg/vk"

	"github.com/urfave/cli"
)

// Profile is command what get full vk user profile
var Profile = cli.Command{
	Name:   "profile",
	Action: runProfile,
}

func runProfile(ctx *cli.Context) error {
	err := setting.NewContext()
	if err != nil {
		log.Printf("err init setting: %v", err)
		return err
	}
	id, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		log.Printf("err atoi %v", err)
		return err
	}

	user, err := vk.GetUser(setting.App.Vk.ServiceToken, int64(id))
	if err != nil {
		log.Printf("err users.get: %v", err)
		return err
	}

	log.Printf("Имя: %s %s", user.FirstName, user.LastName)
	log.Printf("Друзей: %d", user.Counters.Friends)
	log.Printf("Подписчиков: %d", user.FollowersCount)

	t, err := getRegDate(int64(id))
	if err != nil {
		log.Fatalf("Произошла ошибка во премя получения даты регистрации: %v", err)
	}

	log.Printf("Дата регистрации: %s", t.Format("02.01.2006"))

	return nil
}

func getRegDate(id int64) (t time.Time, err error) {
	resp, err := http.Get("https://vk.com/foaf.php?id=" + strconv.Itoa(int(id)))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	log.Printf("%s", data)
	var res Data
	err = xml.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	return res.Person.Created.Date, nil
}

type Data struct {
	XMLName xml.Name `xml:"RDF"`
	Person  struct {
		Created struct {
			Date time.Time `xml:"date,attr"`
		} `xml:"created"`
		Test1 string `xml:"test1"`
	} `xml:"Person"`
}
