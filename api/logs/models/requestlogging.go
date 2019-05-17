package models

type RequestLogging struct {
	Method      string
	Uri         string
	RequestBody string
	Host        string
	UserAgent   string
	RemoteAddr  string
	RemoteIp    string
	Headers     string
	QueryArgs   string
	DateTime    string
	Response    string
	ElapsedTime int64
	StatusCode  int
	ServiceName string
	IndexName   string
}
