package membertype

type MemberLoginReq struct {
	Username string
	Passwd   string
}

type MemberLoginRes struct {
	UserId int64 `json:"userId"`
}

type MemberEntity struct {
	ID       int
	NickName string
}
