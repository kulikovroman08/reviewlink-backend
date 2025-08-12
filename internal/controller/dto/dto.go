package dto

// Структура входного запроса
type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	UserID   string  `json:"-"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Points int    `json:"points"`
}

type CreatePlaceRequest struct {
	Name    string `json:"name" binding:"required"`
	Address string `json:"address" binding:"required"`
}

type CreatePlaceResponse struct {
	ID string `json:"id"`
}

type SubmitReviewRequest struct {
	Token   string `json:"token" binding:"required"`
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Content string `json:"content"`
}
