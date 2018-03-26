// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"
)

type dialogModel int

// Dialogs top level dialogs api
var Dialogs dialogModel

// Get returns dialog by ID
func (dialogModel) Get(dialogID string) (*Dialog, error) {
	return nil, nil
}

type dbDialog struct {
	ID          string `xorm:"unique 'id'"`
	Peer1       int64  `xorm:"unique(dialog_idx)"`
	Peer2       int64  `xorm:"unique(dialog_idx)"`
	Peer1Cursor int64  `xorm:"peer1_cursor"`
	Peer2Cursor int64  `xorm:"peer2_cursor"`
}

type viewDialog struct {
	ID        string `xorm:"unique"`
	Peer1     int64  `xorm:"unique(dialog_idx)"`
	Peer2     int64  `xorm:"unique(dialog_idx)"`
	Text      string
	SenderID  int64                `xorm:"sender_id"`
	CreatedAt *timestamp.Timestamp `xorm:"created_at"`
	MessageID int64                `xorm:"message_id"`
}

func (dbDialog) TableName() string {
	return "dialogs"
}

func (dd viewDialog) Parcipant(userID int64) int64 {
	if dd.Peer1 == userID {
		return dd.Peer2
	}
	return dd.Peer1
}

func chatKey(peer1, peer2 int64) string {
	if peer1 > peer2 {
		peer1, peer2 = peer2, peer1
	}
	return fmt.Sprintf("%d:%d", peer1, peer2)
}

func sortKey(peer1, peer2 int64) (int64, int64) {
	if peer1 > peer2 {
		return peer2, peer1
	}
	return peer1, peer2
}

func (dialogModel) Create(peer1, peer2 int64) error {
	peer1, peer2 = sortKey(peer1, peer2)
	dd := &dbDialog{
		ID:          chatKey(peer1, peer2),
		Peer1:       peer1,
		Peer2:       peer2,
		Peer1Cursor: 0,
		Peer2Cursor: 0,
	}
	_, err := db.InsertOne(dd)
	if err != nil {
		return err
	}
	return nil
}

func (dialogModel) Dialogs(userID int64) (dialogs []*Dialog, err error) {
	var dbDialogs []viewDialog
	sql := `SELECT
  dialogs.id,
  peer1,
  peer2,
  text,
  sender_id,
  created_at,
  messages.id as message_id
FROM dialogs
LEFT JOIN messages
  ON messages.dialog_id = dialogs.id
WHERE peer1 = $1
   OR peer2 = $1
ORDER BY
  messages.id desc
`
	err = db.SQL(sql, userID).Find(&dbDialogs)
	if err != nil {
		return
	}

	return viewDialogsToDialogs(userID, dbDialogs), nil
}

func viewDialogsToDialogs(userID int64, viewDialogs []viewDialog) []*Dialog {
	res := make([]*Dialog, len(viewDialogs))
	for i, dbDialog := range viewDialogs {
		d := &Dialog{
			ParcipantId: dbDialog.Parcipant(userID),
			LastMessage: &Message{
				DialogId:  dbDialog.ID,
				Text:      dbDialog.Text,
				SenderId:  dbDialog.SenderID,
				CreatedAt: dbDialog.CreatedAt,
			},
			//	Readed: viewDialogs[i].Readed,
		}
		res[i] = d
	}
	return res
}

func dbDialogToDialog(dd dbDialog) Dialog {
	return Dialog{}
}

func getPeers(dialogID string) (int64, int64) {
	arr := strings.SplitN(dialogID, ":", 2)
	if len(arr) != 2 {
		return 0, 0
	}
	peer1, _ := strconv.Atoi(arr[0])
	peer2, _ := strconv.Atoi(arr[1])
	return int64(peer1), int64(peer2)
}

func (dialogModel) SetReaded(dialogID string, userID, messageID int64) error {
	sql := `
    UPDATE
      dialogs
    SET
      peer1_cursor = ?
    WHERE peer1 = ?
      AND id = ?`
	_, err := db.Exec(sql, messageID, userID, dialogID)
	if err != nil {
		return err
	}

	sql = `
    UPDATE
      dialogs
    SET
      peer2_cursor = ?
    WHERE peer2 = ?
      AND id = ?`
	_, err = db.Exec(sql, messageID, userID, dialogID)
	if err != nil {
		return err
	}
	return nil
}
