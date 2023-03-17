package user

type User struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Friend []int  `json:"friends"`
}
