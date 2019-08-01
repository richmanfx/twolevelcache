package main

import (
	log "github.com/Sirupsen/logrus"
	"testing"
)

/*
 * Без кеширования - размер кеша равен нулю
 */
func TestZeroCacheSize(t *testing.T) {

	// Выставить параметры логирования (DebugLevel, InfoLevel ...)
	SetLog(log.InfoLevel)

	// Подготовить тестовые данные
	dataSize := 10 // Размер данных
	testData := make([]SimpleStructure, 0, dataSize)
	dataFill(&testData)
	log.Infof("Данные: %+v", testData)

	// Инициализировать кеш с нулевым размером
	cache := CreateSpecifySizeRamCache(0)
	log.Infof("Кеш: %+v", cache)

	// TODO: Пока оставил тестирование - не работает дебаг в Goland-е, отдельным приложением тестирую

}
