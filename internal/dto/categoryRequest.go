package dto

type CreateCategoryRequest struct {
	Name         string `json:"name"`
	ParendId     uint   `json:"parent_id"`
	ImageUrl     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}
