// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"log"
	"strconv"

	"github.com/smolgu/miss/pkg/setting"
	"github.com/smolgu/miss/pkg/vk"

	"github.com/urfave/cli"
)

// Photos command what get top vk user profile
var Photos = cli.Command{
	Name:   "photos",
	Action: runPhotos,
}

func runPhotos(ctx *cli.Context) error {
	err := setting.NewContext()
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(ctx.Args().First())
	if err != nil {
		log.Printf("err atoi %v", err)
		return err
	}
	log.Printf("get photos for %d", id)

	photos, err := vk.GetPhotos(setting.App.Vk.ServiceToken, id)
	if err != nil {
		log.Printf("get photos %v", err)
		return err
	}
	log.Printf("get %d photos", len(photos))
	for _, url := range photos {
		log.Printf("%v", url)
	}
	return nil
}
