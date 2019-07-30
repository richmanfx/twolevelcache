package main

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"time"
)

/* Простая структура */
type SimpleStructure struct {
	Id     int
	Weight float32
	Name   string `faker:"first_name"`
	Player bool
}

func main() {

	// Выставить параметры логирования (DebugLevel, InfoLevel ...)
	//SetLog(log.DebugLevel)
	SetLog(log.InfoLevel)

	// Инициализировать drive-кеш заданного размера
	cacheSize := 1001
	driveCache := CreateSpecifySizeDriveCache(cacheSize)

	// Сгенерировать фейковые данные
	dataAmount := 1000                                   // Количество данных
	cachedData := make([]SimpleStructure, 0, dataAmount) // Слайс для данных
	dataFill(&cachedData)                                // Заполнение данными
	log.Debugf("Данные: %+v", cachedData)

	// Запросить рандомные данные заданное количество раз с использование drive-кеша
	requestAmount := 5000                     // Количество запросов
	rand.Seed(time.Now().Unix())              // Инициализация псевдогенератора временем
	graphicalAnalysisData := make([]int64, 0) // Для сбора данных для графического анализа

	for i := 0; i < requestAmount; i++ {

		// Случайные данные
		randomIndex := rand.Int() % len(cachedData)
		data := cachedData[randomIndex]

		// Получить данные, засечь время получения
		startTime := time.Now().UnixNano()
		findings := getData(driveCache, data)
		finishTime := time.Now().UnixNano()
		receiptTime := finishTime - startTime
		log.Debugf("Полученные данные: %v за время(наносекунды): '%v'", findings, receiptTime)
		graphicalAnalysisData = append(graphicalAnalysisData, receiptTime) // Добавить для графического анализа
	}

	// Вывести график задержек в файл
	dataPlotting(graphicalAnalysisData, cacheSize, dataAmount, requestAmount)

	//// Инициализировать ram-кеш заданного размера
	//cacheSize := 95
	//driveCache := CreateSpecifySizeMemoryCache(cacheSize)
	//
	//// Запросить рандомные данные заданное количество раз с использование кеша
	//requestAmount := 1000                     // Количество запросов
	//rand.Seed(time.Now().Unix())              // Инициализация псевдогенератора временем
	//graphicalAnalysisData := make([]int64, 0) // Для сбора данных для графического анализа
	//
	//for i := 0; i < requestAmount; i++ {
	//
	//	// Случайные данные
	//	randomIndex := rand.Int() % len(cachedData)
	//	data := cachedData[randomIndex]
	//
	//	// Получить данные, засечь время получения
	//	startTime := time.Now().UnixNano()
	//	findings := getData(driveCache, data)
	//	finishTime := time.Now().UnixNano()
	//	receiptTime := finishTime - startTime
	//	log.Infof("Полученные данные: %v за время(наносекунды): '%v'", findings, receiptTime)
	//	graphicalAnalysisData = append(graphicalAnalysisData, receiptTime) // Добавить для графического анализа
	//}
	//
	//// Вывести график задержек в файл
	//dataPlotting(graphicalAnalysisData, cacheSize, dataAmount, requestAmount)

}
