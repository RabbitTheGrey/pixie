package datamapper

import (
	"database/sql"
	"reflect"
)

const TagColumn = "column"

type InvalidDestinationError struct {
	Type reflect.Type
}

func (err *InvalidDestinationError) Error() string {
	return "Invalid destination type."
}

// Сканирует одиночное значение в переданный указатель
//
//	`dest` - указатель на срез &[]T
func SingleScalarResult(row *sql.Row, dest any) error {
	return row.Scan(dest)
}

// Сканирует значения единственного столбца в переданный массив
func SingleColumnResult(rows *sql.Rows, dest []any) error {
	if rows == nil {
		return sql.ErrNoRows
	}
	defer rows.Close()

	slicePointerVal, elemType, err := validateAndExtractType(dest)
	if err != nil {
		return err
	}

	for rows.Next() {
		scanValue := reflect.New(elemType).Interface()
		err = rows.Scan(scanValue)
		if err != nil {
			return err
		}

		slicePointerVal.Set(reflect.Append(slicePointerVal, reflect.ValueOf(scanValue).Elem()))
	}

	return rows.Err()
}

// Сканирует единственную строку в переданную структуру
func SingleResult(row *sql.Row, columns []string, dest any) error {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Pointer || value.IsNil() {
		return &InvalidDestinationError{Type: reflect.TypeOf(dest)}
	}
	value = value.Elem()
	if value.Kind() != reflect.Struct {
		return &InvalidDestinationError{Type: reflect.TypeOf(dest)}
	}

	reflectType := value.Type()
	values := make([]any, len(columns))

	for i, col := range columns {
		fieldPtr, found := matchColumn(value, reflectType, col)
		if found {
			values[i] = fieldPtr
		} else {
			values[i] = new(any) // игнорируем
		}
	}

	return row.Scan(values...)
}

// Сканирует выдачу в переданный массив структур
//
//	`dest` - указатель на слайс структур &[]Struct{} или &[]*Struct{}
func Result(rows *sql.Rows, dest any) error {
	if rows == nil {
		return sql.ErrNoRows
	}
	defer rows.Close()

	slicePointerVal, elemType, isPointer, err := validateAndExtractTypes(dest)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	for rows.Next() {
		structVal, err := createAndScanStruct(rows, columns, elemType)
		if err != nil {
			return err
		}

		appendToSlice(slicePointerVal, structVal, isPointer)
	}

	return rows.Err()
}

// Ищет поле структуры по тэгу
//
//	`column`
//
// и возвращает указатель на него
func matchColumn(value reflect.Value, reflectType reflect.Type, columnName string) (any, bool) {
	for i := 0; i < reflectType.NumField(); i++ {
		field := reflectType.Field(i)
		colName := field.Tag.Get(TagColumn)

		if colName == columnName && value.Field(i).CanSet() {
			return value.Field(i).Addr().Interface(), true
		}
	}

	return nil, false
}

// Создание экземпляра структуры и заполнение его данными из строки
func createAndScanStruct(rows *sql.Rows, columns []string, structType reflect.Type) (reflect.Value, error) {

	structVal := reflect.New(structType).Elem()
	scanValues := make([]any, len(columns))

	for i, colName := range columns {
		column, found := matchColumn(structVal, structType, colName)
		if found {
			scanValues[i] = column
		} else {
			scanValues[i] = new(any)
		}
	}

	err := rows.Scan(scanValues...)
	return structVal, err
}

// Добавление элемента в целевой слайс с учетом типа элемента
func appendToSlice(sliceVal reflect.Value, item reflect.Value, isPointer bool) {
	if isPointer {
		sliceVal.Set(reflect.Append(sliceVal, item.Addr()))
	} else {
		sliceVal.Set(reflect.Append(sliceVal, item))
	}
}

// Проверка dest и получение метаданных о типе элемента среза
func validateAndExtractTypes(dest any) (reflect.Value, reflect.Type, bool, error) {
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return reflect.Value{}, nil, false, &InvalidDestinationError{Type: reflect.TypeOf(dest)}
	}

	value = value.Elem()
	if value.Kind() != reflect.Slice {
		return reflect.Value{}, nil, false, &InvalidDestinationError{Type: reflect.TypeOf(dest)}
	}

	elemType := value.Type().Elem()
	isPointer := elemType.Kind() == reflect.Ptr

	structType := elemType
	if isPointer {
		structType = elemType.Elem()
	}

	if structType.Kind() != reflect.Struct {
		return reflect.Value{}, nil, false, &InvalidDestinationError{Type: reflect.TypeOf(dest)}
	}

	return value, structType, isPointer, nil
}

// Проверка dest и получения типа элемента слайса
func validateAndExtractType(dest any) (reflect.Value, reflect.Type, error) {
	value := reflect.ValueOf(dest)

	if value.Kind() != reflect.Ptr || value.IsNil() {
		return reflect.Value{}, nil, &InvalidDestinationError{
			Type: reflect.TypeOf(dest),
		}
	}

	value = value.Elem()

	if value.Kind() != reflect.Slice {
		return reflect.Value{}, nil, &InvalidDestinationError{
			Type: reflect.TypeOf(dest),
		}
	}

	elemType := value.Type().Elem()
	return value, elemType, nil
}
