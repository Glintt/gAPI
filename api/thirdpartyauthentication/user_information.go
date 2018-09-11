package thirdpartyauthentication

type UserTokenInformation struct {
	Source      string
	Name        string
	Active      bool
	Information string
}

type ProtectionInfo struct {
	Header   string
	UserInfo string
	Cached   bool
	Error    error
}
