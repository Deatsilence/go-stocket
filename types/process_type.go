package types

type ProcessTypes int

const (
	Add ProcessTypes = iota
	Update
	Delete
)
