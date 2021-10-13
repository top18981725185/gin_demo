package vo

type CreategoryRequest struct {
	Name string `json:"name" binding:"required"`
}
