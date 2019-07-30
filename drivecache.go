/* Реализация кеша на диске */

package main

import (
	"encoding/gob"
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
)

const (
	CacheDir = "drive_cache"
)

/* Контейнер */
type DriveCache struct {
	fileNames []string // Имена файла кеша
	maxSize   int      // Максимальный размер кеша
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

	return &DriveCache{fileNames: make([]string, 0), maxSize: size}
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
func (dc *DriveCache) Put(keyFileName string, value interface{}) error {

	// Проверить не заполнен ли кеш полностью
	if dc.maxSize != -1 { // "-1" - нет ограничения в размере кеша
		log.Debugf("Количество элементов в drive-кеше: %d", dc.Size())
		log.Debugf("Максимальный размер drive-кеша: %d", dc.maxSize)

		if dc.Size() >= dc.maxSize {
			log.Infoln("Кеш полностью заполнен - удаляем значение с наименьшей частотой использования!")
			err := dc.LowFrequencyValueDelete()
			if err != nil {
				log.Errorf("Ошибка удаления низкочастотного значения: %s", err)
				return errors.New("ошибка удаления низкочастотного значения")
			}
		}
	}

	// Поместить кешируемый элемент в drive-кеш
	element := &DriveElement{
		Value:     value,
		Frequency: 1, // Помещаем в кеш - значит используется в первый раз
	}

	// Сериалиазовать елемент в файл
	gob.Register(SimpleStructure{}) // Регистрация типа

	fullPath := CacheDir + "/" + keyFileName
	file, err := os.Create(fullPath)
	if err != nil {
		log.Errorf("Ошибка создания файла кеширования '%s': %s", fullPath, err)
		return errors.New("ошибка создания файла кеширования")
	}

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(element)
	if err != nil {
		log.Errorf("Ошибка кодирования: %s", err)
		return errors.New("ошибка 'mob' кодирования")
	}
	err = file.Close()
	if err != nil {
		log.Errorf("Ошибка закрытия файла кеширования '%s': %s", fullPath, err)
		return errors.New("ошибка закрытия файла кеширования")
	}

	// Успешно сериализовали
	dc.fileNames = append(dc.fileNames, keyFileName)
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

	//_, ok := dc.elements[key]
	//if ok {
	//	result = true
	//} else {
	//	result = false
	//}
	//
	return result
}

/* Size */
func (dc *DriveCache) Size() int {
	result := len(dc.fileNames)
	log.Debugf("Количество элементов в drive-кеше: %d", result)
	return result
}

func (dc *DriveCache) LowFrequencyValueDelete() error {

	// Найти самый редкий	//TODO

	// Удалить по ключу

	return nil
}
