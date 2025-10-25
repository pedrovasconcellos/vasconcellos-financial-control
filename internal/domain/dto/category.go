package dto

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=income expense"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
}

type CategoryResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	ParentID    *string `json:"parentId"`
}
