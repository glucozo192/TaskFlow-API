package models

type AddRoleRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateRoleRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}
