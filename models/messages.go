// Copyright 2018 Kirill Zhuharev. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

// TableName return table name in db. Just implement an xorm interface
func (m Message) TableName() string {
	return "messages"
}
