/* Реализация кеша на диске */

package main

import (
	"encoding/gob"
	log "github.com/Sirupsen/logrus"
	"os"
	"sync"
)

const (
	CacheDir = "drive_cache"
)

/* Контейнер */
type DriveCache struct {
	sync.RWMutex                          // Защитный мьютекс для потокобезопасности map
	elements     map[string]*DriveElement // Кешируемые элементы
	maxSize      int                      // Максимальный размер кеша
}

/* Кешируемый элемент */
type DriveElement struct {
	Value     interface{} // Кешируемое значение
	Frequency int         // Частота использования элемента
}

/* Создать новый дисковый кеш заданного размера */
func CreateSpecifySizeDriveCache(size int) Cache {

	// Создать директорию для кеш-файлов, если её нет
	makeDirectory(CacheDir)

	return &DriveCache{elements: make(map[string]*DriveElement), maxSize: size}
}

/* Создать директорию, если её нет */
func makeDirectory(dirName string) {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err = os.Mkdir(dirName, os.ModeDir|0755)
		if err != nil {
			log.Errorf("Не удалось создать директорию '%s': %s", dirName, err)
			panic(err)
		}
	}
}

/* Реализация методов интерфейса Cache */

/* Put */
func (dc *DriveCache) Put(key string, value interface{}) error {

	// Блокировать на время записи
	dc.Lock()
	defer dc.Unlock()

	// Проверить не заполнен ли кеш полностью
	if dc.maxSize != -1 { // "-1" - нет ограничения в размере кеша
		log.Debugf("Количество элементов в кеше: %d", dc.Size())
		log.Debugf("Максимальный размер кеша: %d", dc.maxSize)

		if dc.Size() >= dc.maxSize {
			log.Infoln("Кеш полностью заполнен - удаляем значение с наименьшей частотой использования!")
			err := dc.LowFrequencyValueDelete()
			if err != nil {
				log.Errorf("Ошибка удаления низкочастотного значения: %s", err)
			}
		}
	}

	// Поместить в drive-кеш
	dc.elements[key] = &DriveElement{
		Value:     value,
		Frequency: 1, // Помещаем в кеш - значит используется в первый раз
	}

	// Сериалиазовать структуру в файл
	gob.Register(SimpleStructure{}) // Регистрация типа SimpleStructure
	fileName := key
	serializeValue := dc.elements[key]

	fullPath := CacheDir + "/" + fileName
	file, err := os.Create(fullPath)
	if err != nil {
		log.Errorf("Ошибка создания файла кеширования '%s': %s", fullPath, err)
	}

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(serializeValue)
	if err != nil {
		log.Errorf("Ошибка кодирования: %s", err)
	}
	err = file.Close()
	if err != nil {
		log.Errorf("Ошибка закрытия файла кеширования '%s': %s", fullPath, err)
	}

	return nil
}

/* Get */
func (dc *DriveCache) Get(key string) interface{} {
	var result interface{}
	return result
}

/* Del */
func (dc *DriveCache) Del(key string) error {
	var result error
	return result
}

/* IsExist */
func (dc *DriveCache) IsExist(key string) bool {
	var result bool
	return result
}

/* Size */
func (dc *DriveCache) Size() int {
	result := len(dc.elements)
	log.Debugf("Количество элементов в drive-кеше: %d", result)
	return result
}

func (dc *DriveCache) LowFrequencyValueDelete() error {

	// Найти самый редкий	//TODO

	// Удалить по ключу

	return nil
}
