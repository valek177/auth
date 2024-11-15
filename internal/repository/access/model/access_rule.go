package model

// AccessRule is a struct for access rule record
type AccessRule struct {
	Role     string `json:"role"`
	Endpoint string `json:"endpoint"`
}
