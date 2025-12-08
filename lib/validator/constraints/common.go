package constraints

type Constraint struct {
	// Общие правила валидации
	Type     *string
	Nullable *bool
}

// Добавление валидации типа
//
//	Переменная `type` недоступна из-за совпадения с ключевым словом, поэтому будет просто `val`
func (constraint *Constraint) WithType(val string) *Constraint {
	constraint.Type = &val
	return constraint
}

// Добавление валидации пустого значения
func (constraint *Constraint) WithNullable(nullable bool) *Constraint {
	constraint.Nullable = &nullable
	return constraint
}
