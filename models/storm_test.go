package models_test

import (
	"os"
	"testing"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
)

type TestStruct struct {
	ID int64 `json:"id" storm:"index"`
}

func TestNotIn(t *testing.T) {
	t.Run("Open", func(t *testing.T) {
		db, err := storm.Open("/tmp/test.storm")
		if err != nil {
			t.Fatal(err)
		}

		for i := 1; i < 16; i++ {
			ts := TestStruct{
				ID: int64(i),
			}
			err = db.Save(&ts)
			if err != nil {
				t.Fatal(err)
			}
		}

		var res []TestStruct
		err = db.Select(q.Not(q.In("ID", []int64{4, 6, 3}))).Find(&res)
		if err != nil {
			t.Fatal(err)
		}
		if len(res) != 12 {
			t.Fatalf("len res not 12, got %d", len(res))
		}
		for _, v := range res {
			if v.ID == 4 || v.ID == 6 || v.ID == 3 || v.ID == 0 {
				t.Fatalf("id = %d", v.ID)
			}
		}

		defer os.Remove("/tmp/test.storm")
	})

}
