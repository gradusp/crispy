package model

type SecurityZone struct {
	tableName struct{}   `pg:"controller.security_zones"`
	ID        string     `json:"id" pg:"id,pk,type:uuid,default:gen_random_uuid()"`
	Name      string     `json:"name" pg:"name,unique"`
	Clusters  []*Cluster `pg:"rel:has-many"`
}
