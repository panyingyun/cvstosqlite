package model

type Node struct {
	Id   int    `json:"id" xorm:"not null pk autoincr INT(11)"`
	RFU1 string `json:"rfu1" xorm:"VARCHAR(100) 'rfu1'"`
	RFU2 string `json:"rfu2" xorm:"VARCHAR(100) 'rfu2'"`
	RFU3 string `json:"rfu3" xorm:"VARCHAR(100) 'rfu3'"`
}
