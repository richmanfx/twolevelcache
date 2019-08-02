package main

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

/* Расчёт хеша структуры */
func TestGetHash(t *testing.T) {

	// Выставить параметры логирования (DebugLevel, InfoLevel ...)
	setLog(log.InfoLevel)

	// Подготовить тестовые данные
	type TestData struct {
		Id             int
		Weight         float32
		Name           string
		Player         bool
		expectedResult string
	}
	var testDataSlice = make([]TestData, 0)

	testDataSlice = append(testDataSlice,
		TestData{Id: 0, Weight: 0.0, Name: "Maya", Player: true, expectedResult: "15543582922617197410"},
		TestData{Id: 1, Weight: 1.1, Name: "Selina", Player: false, expectedResult: "9603443469398688814"},
		TestData{Id: 2, Weight: 12.12, Name: "Estrella", Player: true, expectedResult: "4044771355162859765"},
		TestData{Id: 3, Weight: 123.12345678, Name: "Hermina", Player: false, expectedResult: "7988528797218559496"},
	)

	// Пробежать по всем данным
	for _, testDataItem := range testDataSlice {

		simpleStructure := SimpleStructure{
			Id:     testDataItem.Id,
			Weight: testDataItem.Weight,
			Name:   testDataItem.Name,
			Player: testDataItem.Player,
		}

		// Фактический результат
		actualResult := getHash(simpleStructure)

		// Сравнение ожидаемого и фактического результатов
		if actualResult != testDataItem.expectedResult {
			t.Errorf("Неверно вычисленный ХЕШ '%s'. Ожидается '%s'.", actualResult, testDataItem.expectedResult)
		} else {
			t.Logf("Верно вычисленный ХЕШ '%s'. Ожидался '%s'.", actualResult, testDataItem.expectedResult)
		}

	}

}
