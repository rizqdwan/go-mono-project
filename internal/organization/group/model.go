package group

type Group struct {
	ID        int64  `db:"id"`
	Label     string `db:"label"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type GroupResponse struct {
	ID        int64  `json:"id"`
	Label     string `json:"label"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
