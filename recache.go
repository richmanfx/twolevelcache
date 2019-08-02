package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

/* Рекеширование */
func reCaching(mc *MemoryCache) error {

	dc := &DriveCache{fileNames: nil, maxSize: driveCacheSize}

	// Перенести в DRIVE-кеш если частота меньше или равна Минимальной частоте нахождения элемента в RAM-кеше
	for _, element := range mc.elements {
		if element.Frequency <= minRamFrequency {
			key := getHash(element.Value.(SimpleStructure))

			// Внести в DRIVE-кеш
			err := dc.Put(key, element.Value)
			if err != nil {
				log.Errorf("Ошибка внесения в DRIVE-кеш при рекешинге : %s", err)
				return errors.New("ошибка внесения в DRIVE-кеш при рекешинге")
			}

			// Удалить из RAM-кеша
			err = mc.Del(key)
			if err != nil {
				log.Errorf("Ошибка удаления из RAM-кеша при рекешинге : %s", err)
				return errors.New("ошибка удаления из RAM-кеша при рекешинге")
			}
		}
	}

	// Перенести из DRIVE-кеша в RAM-кеш если частота больше Минимальной частоты нахождения элемента в RAM-кеше
	for _, key := range dc.getAllKeys() { // Пробежаться по DRIVE-кешу

		element := dc.Get(*key)

		if element.Frequency > minRamFrequency {
			// Внести в RAM-кеш
			err := mc.Put(*key, element)
			if err != nil {
				log.Errorf("Ошибка внесения в RAM-кеш при рекешинге : %s", err)
				return errors.New("ошибка внесения в RAM-кеш при рекешинге")
			}

			// Удалить из DRIVE-кеша
			err = dc.Del(key)
			if err != nil {
				log.Errorf("Ошибка удаления из DRIVE-кеша при рекешинге : %s", err)
				return errors.New("ошибка удаления из DRIVE-кеша при рекешинге")
			}
		}
	}

	return nil
}
