package usertype

type UserLoginReq struct {
	Username string
	Passwd   string
}

type UserLoginRes struct {
	UserId int64 `json:"userId"`
}
