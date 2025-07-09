package types

type Request struct {
	Email string `json:"email" validator:"required,email"`
	Hash  string `json:"hash,omitempty"`
}
