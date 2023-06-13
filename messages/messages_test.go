package messages

import (
	"fmt"
	"testing"
	"time"
)

func Test_Message(t *testing.T) {

	/* 	MS().PutMessage("Hello", NewMessage())
	   	msgs, err := MS().GetMessage("Hello", "你有啥绝活")
	   	if err != nil {
	   		t.Errorf("error of %s", err)
	   	}
	   	fmt.Println(msgs) */
	/* 	t2 := time.Now().Add(time.Second)
	   	time.Sleep(time.Second)
	   	fmt.Println(t2.After(time.Now())) */

	t3, err := time.Parse("2006-01-02 15:04:05", "2023-06-13 14:33:00")
	if err != nil {
		fmt.Println(err)
		return
	}
	t4, err := time.Parse("2006-01-02 15:04:05", "2023-06-13 14:33:14")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(t3.Before(t4))

}
