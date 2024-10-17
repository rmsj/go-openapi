package models

// Topic of a thread.
type Topic struct {
	Namespace string `json:"namespace" validate:"required"`
	Topic     string `json:"topic"  validate:"omitempty,oneof=general politics health"`
	Private   bool   `json:"private"`
	ViewCount int64  `json:"viewCount"`
}
