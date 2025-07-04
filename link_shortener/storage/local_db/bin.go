package local_db

type Bin struct {
	Email string `json:"email" validate:"required,email"`
	Hash  string `json:"hash" validate:"required"`
}

func NewBin(email string, hash string) *Bin {
	return &Bin{
		Email: email,
		Hash:  hash,
	}
}
