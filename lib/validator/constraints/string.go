package constraints

type String struct {
	Constraint

	// Правила валидации строковых переменных
	MinLen *int
	MaxLen *int
}

// Добавление валидации минимальной длины строки
func (constraint *String) WithMinLen(minLen int) *String {
	constraint.MinLen = &minLen
	return constraint
}

// Добавление валидации максимальной длины строки
func (constraint *String) WithMaxLen(maxLen int) *String {
	constraint.MaxLen = &maxLen
	return constraint
}
