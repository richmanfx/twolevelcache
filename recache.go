package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
)

func reCaching(mc *MemoryCache) error {

	// Минимальная частота нахождения элемента в RAM-кеше:
	//  - если меньше или равно, то переносить в Drive-кеш при рекешировании
	//  - если больше, то переносить из Drive-кеша в Ram-кеш

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

	//// Перенести из Drive-кеша в Ram-кеш если частота больше Минимальной частоте нахождения элемента в RAM-кеше

	// Получить все ключи из DRIVER-кеш
	//dc.getAllKey()

	return nil
}
