package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bxcodec/faker"
	"github.com/mitchellh/hashstructure"
	. "os"
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

	// Заполнить слайс данными
	dataSize := 10
	cachedData := make([]SimpleStructure, 0, dataSize)
	dataFill(&cachedData)
	log.Debugf("Данные: %+v", cachedData)

	// Инициализировать кеш с нулевым размером
	//cache := CreateMemoryCache()
	cache := CreateSpecifySizeMemoryCache(0)

	log.Infof("Кеш: %+v", cache)

	// Закешировать значение
	data := cachedData[0]
	hash := getHash(data)
	log.Infof("Хеш: %+v", hash)
	err := cache.Put(hash, data)
	if err != nil {
		log.Infof("Ошибка добавления в кеш: %s", err)
	}
	log.Infof("Кеш: %+v", cache)

	// Закешировать значение
	data = cachedData[1]
	hash = getHash(data)
	log.Infof("Хеш: %+v", hash)
	err = cache.Put(hash, data)
	if err != nil {
		log.Infof("Ошибка добавления в кеш: %s", err)
	}
	log.Infof("Кеш: %+v", cache)

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
