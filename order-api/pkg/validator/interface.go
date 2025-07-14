package validator

type Validatable interface {
	Validate() error
}

func ValidateModel(model Validatable) error {
	return model.Validate()
}
