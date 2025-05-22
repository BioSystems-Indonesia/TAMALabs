package entity

type Session struct {
	ID        int64  `json:"id"`
	AdminID   int64  `json:"admin_id"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
