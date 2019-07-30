/* Реализация кеша на диске */

package main

import (
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

	// Кешируемый элемент
	element := &MemoryElement{
		Value:     value,
		Frequency: 1, // Помещаем в кеш - значит используется в первый раз
	}

	// Сериалиазовать элемент в файл
	err := gobEncode(keyFileName, element)
	if err != nil {
		log.Errorf("Ошибка сериализации файла: %s", err)
		panic(err)
	}

	// Успешно сериализовали
	dc.fileNames = append(dc.fileNames, keyFileName)
	return nil
}

/* Get */
func (dc *DriveCache) Get(key string) interface{} {
	var result interface{}

	element, err := gobDecode(key)
	if err == nil {
		log.Debugf("Получен из drive-кеша элемент '%v'", element)
		element.Frequency++ // Частота использования
		result = element.Value
	} else {
		log.Errorf("Ошибка десериализации файла: %s", err)
		panic(err)
	}

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

	for _, fileName := range dc.fileNames {
		if key == fileName {
			result = true
			log.Debugf("Элемент '%s' уже находится в drive-кеше", key)
		}
	}

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
