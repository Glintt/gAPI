package models

import "encoding/json"

type RequestLogging struct {
	Id          string
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
	JWTContent  string
}

func (r *RequestLogging) GetOtherInfo() string {
	var info = make(map[string]string)

	info["jwt_content"] = r.JWTContent
	i, err := json.Marshal(info)
	if err != nil {
		return "{}"
	}
	return string(i)
}
