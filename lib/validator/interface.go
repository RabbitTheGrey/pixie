package validator

type IConstraint interface {
	validate()
	validateString()
}
