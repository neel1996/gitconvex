package validator

type Validator interface {
	Validate(interface{}) error
}

type ValidatorWithStringFields interface {
	ValidateWithFields(fields ...string) error
}
