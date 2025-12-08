package example

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Пример тела запроса
type RequestData struct {
	Something *string `json:"something"`
}

// Набор тестовых данных
func dataFixtures() []string {
	data := []string{
		"foo",
		"bar",
		"baz",
	}

	return data
}

// Получение данных списком
func List(w http.ResponseWriter, r *http.Request, params map[string]string) {
	jsonResponse, err := json.Marshal(dataFixtures())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

// Получение строки из набора данных по индексу
func Get(w http.ResponseWriter, r *http.Request, params map[string]string) {
	strIndex, ok := params["index"]
	if !ok {
		http.Error(w, "Не найден обязательный параметр `index`", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(strIndex)

	if err != nil {
		http.Error(w, "Параметр `index` должен иметь тип `int`", http.StatusBadRequest)
		return
	}

	var data []string = dataFixtures()

	if index >= len(data) {
		http.NotFound(w, r)
		return
	}

	var elem string = data[index]

	jsonResponse, err := json.Marshal(elem)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

// Добавление данных в тестовый массив и получение обновленного списка
func Post(w http.ResponseWriter, r *http.Request, params map[string]string) {
	var body RequestData

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)

	if err != nil {
		http.Error(w, "Не удалось декодировать JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if body.Something == nil {
		http.Error(w, "Не передано обязательное поле something", http.StatusBadRequest)
		return
	}

	data := dataFixtures()
	data = append(data, *body.Something)

	jsonResponse, err := json.Marshal(data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
