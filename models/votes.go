package models

type votesModel int

var Votes votesModel

func (votesModel) Vote(targetID int64, vote VoteType, voterID int64) (matched bool, err error) {
	return
}
