package models

import "time"

type votesModel int

var Votes votesModel

type userVoteModel struct {
	TargetID int64 `xorm:"'target_id' unique(idx_vote)"`
	VoterID  int64 `xorm:"'voter_id' unique(idx_vote)"`
	Vote     VoteType
	VotedAt  time.Time
}

func (userVoteModel) TableName() string {
	return "votes"
}

func (votesModel) Vote(targetID int64, vote VoteType, voterID int64) (matched bool, err error) {
	v := userVoteModel{
		TargetID: targetID,
		VoterID:  voterID,
		VotedAt:  time.Now(),
		Vote:     vote,
	}
	_, err = db.Insert(&v)
	if err != nil {
		return
	}

	cnt, err := db.Where("voter_id = ? AND target_id = ? AND vote = ?", targetID, voterID, VoteType_like).Count()
	if err != nil {
		return
	}

	return cnt == 1, nil
}
