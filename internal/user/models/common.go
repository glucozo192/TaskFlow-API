package models

type LoginRequest struct {
	Username string `db:"username" json:"username"`
	Password string `db:"password" json:"password"`
}

type LoginResponse struct {
	User  User
	Token string `db:"token" json:"token"`
}

type CountResponse struct {
	Count int64 `json:"count,omitempty" db:"count"`
}

type GetListPermissionsRow struct {
	RoleID      string   `db:"role_id" json:"role_id"`
	Permissions []string `db:"permissions" json:"permissions"`
}
