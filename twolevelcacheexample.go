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

const (
	ramCacheSize    = 20  // Размер RAM-части кеша
	driveCacheSize  = 40  // Размер DRIVE-части кеша
	minRamFrequency = 5   // Минимальная частота нахождения элемента в RAM-части кеша
	minFrequency    = 3   // Минимальная частота нахождения элемента кеше (должна быть меньше, чем в RAM-части)
	dataNumber      = 50  // Количество данных
	requestNumber   = 300 // Количество запросов

	fromAfarDelay       = 500                  // Задержка получения данных "издалека" когда нет данных в кеше, в мкс
	graphResultFileName = "requests_delay.jpg" // Результаты в графическом виде
	cacheDir            = "drive_cache"        // Директория для DRIVE-кеша

	// Количество нерезультативных запросов в кеш, после которого запускать рекеширование
	//recacheRequestsNumber = ramCacheSize / 2
	recacheRequestsNumber = 50
)

func main() {

	// Выставить параметры логирования (DebugLevel, InfoLevel ...)
	//SetLog(log.DebugLevel)
	SetLog(log.InfoLevel)

	// Инициализировать кеш заданного размера
	twoLevelCache := CreateTwoLevelCache(ramCacheSize, driveCacheSize, minRamFrequency, minFrequency)

	// Сгенерировать фейковые данные
	cachedData := make([]SimpleStructure, 0, dataNumber) // Слайс для данных
	dataFill(&cachedData)                                // Заполнение данными
	log.Debugf("Данные: %+v", cachedData)

	// Запросить рандомные данные заданное количество раз
	rand.Seed(time.Now().Unix())              // Инициализация псевдогенератора временем
	graphicalAnalysisData := make([]int64, 0) // Для сбора данных для графического анализа

	for i := 0; i < requestNumber; i++ {

		// Случайные данные
		randomIndex := rand.Int() % len(cachedData)
		data := cachedData[randomIndex]

		// Получить данные, засечь время получения
		startTime := time.Now().UnixNano()
		findings := twoLevelCache.Get(data)

		// Если данных в кеше нет
		if findings == nil {

			// Эмуляция получения данных не из кеша - задержка
			time.Sleep(fromAfarDelay * time.Microsecond)
			findings = &MemoryElement{Value: data}

			// Занести в кеш
			err := twoLevelCache.Put(data)
			if err != nil {
				log.Infof("Ошибка помещения в кеш: %s", err)
			}
		}

		finishTime := time.Now().UnixNano()
		receiptTime := finishTime - startTime
		log.Infof("Полученные данные, %d из %d): %v за время (наносекунды): '%v'",
			i+1, requestNumber, findings, receiptTime)

		// Для графического анализа
		graphicalAnalysisData = append(graphicalAnalysisData, receiptTime)

	}

	// Вывести график задержек в файл
	dataPlotting(graphicalAnalysisData, driveCacheSize, dataNumber, requestNumber)

}
