package form

// GeneralGetDto General get dto
type GeneralGetDto struct {
	ID int64 `uri:"id" json:"id" binding:"required"`
}
