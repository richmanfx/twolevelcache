package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bxcodec/faker"
	. "os"
)

/* Простая структура */
type SimpleStructure struct {
	Id		int
	Weight	float32
	Name	string		`faker:"first_name"`
	Player	bool
}

func main() {

	// Выставить параметры логирования (DebugLevel, InfoLevel ...)
	SetLog(log.DebugLevel)

	// Заполнить слайс данными
	dataSize := 5
	cachedData := make([]SimpleStructure, 0, dataSize)
	dataFill(&cachedData)
	log.Debugf("Данные: %+v", cachedData)

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
