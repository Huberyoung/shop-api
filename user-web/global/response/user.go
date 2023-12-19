package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var temp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006年01月02日 15:04:05"))
	return []byte(temp), nil
}

type UserResponse struct {
	Id       int32  `json:"id"`
	Mobile   string `json:"mobile"`
	NickName string `json:"nick_name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Gender   int      `json:"gender"`
}
