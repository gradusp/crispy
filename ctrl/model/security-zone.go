package model

type SecurityZone struct {
	ID   string `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name string `json:"name" pg:"name,unique"`
}
