package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bxcodec/faker"
	"github.com/mitchellh/hashstructure"
	. "os"
	"time"
)

/* Получить данные, при возможности воспользоваться кешем */
func getData(cache *DriveCache, data SimpleStructure) interface{} {

	key := getHash(data)
	if !cache.IsExist(key) {
		// Эмуляция получения данных не из кеша - задержка
		time.Sleep(500 * time.Microsecond)
		err := cache.Put(key, data) // Занести в кеш
		if err != nil {
			log.Infof("Ошибка добавления в кеш: %s", err)
		}
	}
	return cache.Get(key)
}

/* Вернуть хеш структуры */
func getHash(structure SimpleStructure) string {
	hash, err := hashstructure.Hash(structure, nil)
	if err != nil {
		log.Errorf("Ошибка хеширования: %s", err)
	}
	//return fmt.Sprintf("%d%d", hash, rand.Int())		// TODO:
	return fmt.Sprintf("%d", hash) // TODO:
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
