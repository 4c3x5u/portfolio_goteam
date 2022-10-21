package requests

type Register struct {
	Usn string `json:"username"`
	Pwd string `json:"password"`
	Ref string `json:"referrer"`
}
