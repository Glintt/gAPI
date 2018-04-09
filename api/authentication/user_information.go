package authentication

type UserTokenInformation struct {
	Source      string `json:"source"`
	Name        string `json:"name"`
	Active      bool   `json:"active"`
	Information string
}

type ProtectionInfo struct {
	Header   string
	UserInfo string
	Cached bool
	Error    error
}
