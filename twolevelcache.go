package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bxcodec/faker"
	"github.com/mitchellh/hashstructure"
	"math/rand"
	. "os"
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
	SetLog(log.InfoLevel)

	// Сгенерировать фейковые данные
	dataAmount := 100                                    // Количество данных
	cachedData := make([]SimpleStructure, 0, dataAmount) // Слайс для данных
	dataFill(&cachedData)                                // Заполнение данными
	log.Debugf("Данные: %+v", cachedData)

	// Инициализировать drive-кеш заданного размера
	cacheSize := 10
	driveCache := CreateSpecifySizeDriveCache(cacheSize)

	// Запросить рандомные данные заданное количество раз с использование drive-кеша
	requestAmount := 10                       // Количество запросов
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
		log.Infof("Полученные данные: %v за время(наносекунды): '%v'", findings, receiptTime)
		graphicalAnalysisData = append(graphicalAnalysisData, receiptTime) // Добавить для графического анализа
	}

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

/* Получить данные, при возможности воспользоваться кешем */
func getData(cache Cache, data SimpleStructure) interface{} {

	key := getHash(data)
	if cache.IsExist(key) {
		return cache.Get(key)
	} else {
		// Эмуляция получения данных не из кеша - задержка
		time.Sleep(500 * time.Microsecond)
		err := cache.Put(key, data) // Занести в кеш
		if err != nil {
			log.Infof("Ошибка добавления в кеш: %s", err)
		}
		return data
	}
}

/* Вернуть хеш структуры */
func getHash(structure SimpleStructure) string {
	hash, err := hashstructure.Hash(structure, nil)
	if err != nil {
		log.Errorf("Ошибка хеширования: %s", err)
	}
	return fmt.Sprintf("%d", hash)
}

/* Заполнить слайс данными */
func dataFill(allCachedData *[]SimpleStructure) {
	var individualCachedData SimpleStructure
	for i := 0; i < cap(*allCachedData); i++ {
		err := faker.FakeData(&individualCachedData)
		if err != nil {
			log.Errorf("Ошибка генерации фейковых данных: %s", err)
		}
		*allCachedData = append(*allCachedData, individualCachedData)
	}
}

/* Выставить параметры логирования */
func SetLog(debugLevel log.Level) {
	log.SetOutput(Stdout)
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006/01/02 15:04:05"
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(debugLevel)
}
