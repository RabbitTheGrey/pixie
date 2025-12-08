package constraints

type Integer struct {
	Constraint

	// Правила валидации целочисленных переменных
	Range    *[2]int
	Positive *bool
	Max      *int
	Min      *int
}

// Добавление валидации границ числового значения
//
//	Переменная `range` недоступна из-за совпадения с ключевым словом, поэтому будет просто `val`
func (constraint *Integer) WithRangle(val [2]int) *Integer {
	constraint.Range = &val
	return constraint
}

// Добавление валидации положительного числа
func (constraint *Integer) WithPositive(positive bool) *Integer {
	constraint.Positive = &positive
	return constraint
}

// Добавление валидации минимального значения числа
func (constraint *Integer) WithMin(min int) *Integer {
	constraint.Min = &min
	return constraint
}

// Добавление валидации максимального значения числа
func (constraint *Integer) WithMax(max int) *Integer {
	constraint.Max = &max
	return constraint
}
