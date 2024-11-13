package model

// AccessRule is a struct for access rule record
type AccessRule struct {
	Id       int64  `json:"id"`
	Role     int64  `json:"role"`
	Endpoint string `json:"endpoint"`
}
