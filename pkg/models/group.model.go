package models

type Group struct {
	JsonEmbeddable
	ID              uint32   `json:"id" validate:"-" sql:"id"`
	Type            string   `json:"type" validate:"required" sql:"type"` // investor or startup team
	Name            string   `json:"group_name" validate:"required" sql:"group_name"`
	Owner           uint32   `json:"group_owner" validate:"required"`
	Bio             string   `json:"group_bio" validate:"required" sql:"group_bio"`
	Tags            []string `json:"tags" validate:"required"`
	NumGroupMembers int      `json:"num_group_members" validate:"required" sql:"num_group_members"`
	GroupMembers    []string `json:"group_member_ids" valudate:"-"`
}
