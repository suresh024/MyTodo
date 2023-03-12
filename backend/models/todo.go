package models

type Todo struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Status      string `json:"status" bson:"status"`
	Audit       Audit  `json:"audit" bson:"audit"`
}

type Audit struct {
	CreatedAt int    `json:"created_at" bson:"created_at"`
	CreatedBy string `json:"created_by" bson:"created_by"`
	UpdateAt  int    `json:"update_at" bson:"update_at"`
	UpdateBy  string `json:"update_by" bson:"update_by"`
}
