/* Реализация кеша на диске */

package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
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

	// Очистить drive-кеш
	clearDriveCache(CacheDir)

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

/* Очистить drive-кеш */
func clearDriveCache(dirName string) {
	dir, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Errorf("Не удалось очистить директорию '%s': %s", dirName, err)
		panic(err)
	}

	fileNumber := 0
	for i, d := range dir {
		err := os.RemoveAll(path.Join([]string{CacheDir, d.Name()}...))
		if err != nil {
			log.Errorf("Не удалось удалить файлы из директори '%s': %s", dirName, err)
			panic(err)
		}
		fileNumber += i
	}
	log.Infof("Удалено %d кеш-файлов в директории %s", fileNumber, dirName)

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

		// Инкрементировать частота использования и обновить в кеше
		log.Debugf("Частота до инкремента: %d", element.Frequency)
		element.Frequency++
		log.Debugf("Частота после инкремента: %d", element.Frequency)
		err = dc.Update(key, element)
		if err != nil {
			log.Errorf("Ошибка обновления значения в drive-кеше: %s", err)
			panic(err)
		}
		log.Debugln("Кеш удачно обновлён")

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

/* Update - обновление значения в кеше */
func (dc *DriveCache) Update(fileName string, element interface{}) error {

	err := gobEncode(fileName, element)
	if err != nil {
		log.Errorf("Ошибка сериализации файла при обновлении drive-кеша: %s", err)
		panic(err)
	}
	return nil
}

func (dc *DriveCache) LowFrequencyValueDelete() error {

	// Найти самый редкий	//TODO

	// Удалить по ключу

	return nil
}
