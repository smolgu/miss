package models

import (
	"github.com/asdine/storm/q"
)

type userModel int

var Users userModel

func (userModel) Get(userID int64) (*User, error) {
	return nil, nil
}

func (userModel) GetByVkID(vkID int64) (*User, error) {
	return nil, nil
}

func (userModel) Random(voterID int64) (res []*User, err error) {
	var votes []Vote
	err = stormDB.Select(q.Eq("VoterId", voterID)).Find(&votes)
	if err != nil {
		return
	}
	var ids []int64
	for _, v := range votes {
		ids = append(ids, v.TargetUserId)
	}
	err = stormDB.Select(q.Not(q.In("ID", ids))).Find(&res)
	if err != nil {
		return
	}
	return
}
