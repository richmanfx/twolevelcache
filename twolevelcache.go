/* Реализация двухуровневого кеша */

package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

type TwoLevelCache struct {
	ramCache          *MemoryCache
	driveCache        *DriveCache
	minRamFrequency   int
	minFrequency      int
	badQueriesCounter int
}

/* Создать новый двухуровневый кеш */
func CreateTwoLevelCache(ramCacheSize, driveCacheSize, minRamFrequency, minFrequency int) *TwoLevelCache {

	return &TwoLevelCache{
		ramCache:          CreateSpecifySizeRamCache(ramCacheSize),
		driveCache:        CreateSpecifySizeDriveCache(driveCacheSize),
		minRamFrequency:   minRamFrequency,
		minFrequency:      minFrequency,
		badQueriesCounter: 0, // Счётчик количества запросов данных, которых не было в кеше
	}

}

/* Получить данные из кеша */
func (tlc *TwoLevelCache) Get(data SimpleStructure) *MemoryElement {

	var result *MemoryElement
	key := getHash(data)

	// Поиск в RAM-кеше
	inRamExist := tlc.ramCache.IsExist(key)
	if inRamExist {
		// Читать из RAM-кеша
		result = tlc.ramCache.Get(key)
	} else {
		// Поиск в Drive-кеше
		inDriveExist := tlc.driveCache.IsExist(key)
		if inDriveExist {
			// Читать из DRIVE-кеша
			log.Debugf("Есть в DRIVE-кеше - читаю")
			result = tlc.driveCache.Get(key)
		} else {
			// В кеше нет запрашиваемых данных
			tlc.badQueriesCounter++ // Увеличить счётчик количества запросов данных, которых не было в кеше
			result = nil
		}
	}

	// Рекеширование
	// TODO: В горутине запустить?
	if tlc.badQueriesCounter > recacheRequestsNumber { // Когда данных не нашлось более заданного количества раз
		log.Infof("Данных не нашлось в кеше более '%d' раз - рекешируем", recacheRequestsNumber)
		err := reCaching(tlc.ramCache)
		if err != nil {
			log.Infof("Ошибка рекеширования: %s", err)
		}
	}

	// Выдать данные
	return result
}

/* Внести данные в кеш */
func (tlc *TwoLevelCache) Put(data SimpleStructure) error {
	element := &MemoryElement{Value: data, Frequency: 1} // Не было в кеше - значит используется в первый раз
	key := getHash(data)

	// Внести в кеш - вносим только в RAM-часть
	err := tlc.ramCache.Put(key, element)
	if err != nil {
		log.Infof("Ошибка помещения в RAM-кеш: %s", err)
		return errors.New("ошибка помещения в RAM-кеш")
	}

	return nil
}
