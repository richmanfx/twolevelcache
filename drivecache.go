/* Реализация кеша на диске */

package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"sort"
)

const (
	CacheDir = "drive_cache"
)

/* Контейнер */
type DriveCache struct {
	fileNames []*string // Имена файла кеша
	maxSize   int       // Максимальный размер кеша
}

/* Создать новый дисковый кеш заданного размера */
func CreateSpecifySizeDriveCache(size int) *DriveCache {

	// Создать директорию для кеш-файлов, если её нет
	makeDirectory(CacheDir)

	// Очистить drive-кеш
	clearDriveCache(CacheDir)

	return &DriveCache{fileNames: make([]*string, 0), maxSize: size}
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
	for _, d := range dir {
		err := os.RemoveAll(path.Join([]string{CacheDir, d.Name()}...))
		if err != nil {
			log.Errorf("Не удалось удалить файлы из директории '%s': %s", dirName, err)
			panic(err)
		}
		fileNumber++
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
			log.Debugln("Кеш полностью заполнен - удаляем значение с наименьшей частотой использования!")
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

	// Добавить запись если такой ещё нет
	elementExist := false
	for _, fileName := range dc.fileNames {

		if *fileName == keyFileName {
			elementExist = true
			break
		}
	}
	if !elementExist {
		dc.fileNames = append(dc.fileNames, &keyFileName)
	}

	return nil
}

/* Get */
func (dc *DriveCache) Get(key string) *MemoryElement {

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

		return element
	} else {
		log.Errorf("Ошибка десериализации файла: %s", err)
		panic(err)
	}

}

/* Del */
func (dc *DriveCache) Del(keyFileName *string) error {
	var result error

	// Удалить файл
	err := os.Remove(CacheDir + "/" + *keyFileName)
	if err != nil {
		log.Errorf("Не удалось удалить файл '%s' из директории '%s': %s", *keyFileName, CacheDir, err)
		result = errors.New(
			fmt.Sprintf("ошибка удаления файла '%s' из директории '%s': %s", *keyFileName, CacheDir, err))
	}
	log.Debugf("Файл '%s' удалён", *keyFileName)

	// Удалить из контейнера
	dc.fileNames = dc.remove(dc.fileNames, keyFileName)
	log.Debugf("Запись '%s' из контейнера удалена", *keyFileName)

	return result
}

/* Удалить ключ-имя_файла из массива ссылок  */
func (dc *DriveCache) remove(fileNames []*string, keyFileName *string) []*string {

	for i, fileName := range fileNames {
		if *fileName == *keyFileName {
			copy(fileNames[i:], fileNames[i+1:])

			// Хвост очистить
			fileNames[len(fileNames)-1] = nil
			fileNames = fileNames[:len(fileNames)-1]
			break
		}
	}
	return fileNames
}

/* IsExist */
func (dc *DriveCache) IsExist(key string) bool {
	var result bool

	for _, fileName := range dc.fileNames {
		if &key == fileName {
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

/* Удалить из drive-кеша значение с наименьшей частотой использования */
func (dc *DriveCache) LowFrequencyValueDelete() error {

	// Массив структур с именами файлов и соответствующими частотами
	type frequencyStruct struct {
		fileName  string
		frequency int
	}
	var frequencyArray []frequencyStruct

	// Считать из drive-кеша
	log.Debugf("До удаления ===> dc.fileNames: %v", dc.fileNames)
	for i, fileName := range dc.fileNames {
		log.Debugf("Имя в dc.fileNames До удаления ===> %d: %v", i, *fileName)
	}

	for _, fileName := range dc.fileNames {
		element := dc.Get(*fileName)
		log.Debugf("В методе LowFrequencyValueDelete получен элемент: %v", element)
		frequencyArray = append(frequencyArray, frequencyStruct{*fileName, element.Frequency})
	}

	log.Debugf("структура 'имя_файла:частота': %v", frequencyArray)

	// Отсортировать
	sort.SliceStable(frequencyArray, func(i, j int) bool {
		return frequencyArray[i].frequency < frequencyArray[j].frequency
	})

	log.Debugf("Минимальная частота: %v", frequencyArray[0])

	// Удалить по ключу
	err := dc.Del(&frequencyArray[0].fileName)
	if err != nil {
		log.Errorf("Ошибка удаления элемента из drive-кеша: %s", err)
		panic(err)
	}

	log.Debugf("После удаления ===> dc.fileNames: %v", len(dc.fileNames))
	log.Debugf("Элемент '%s' удалён из drive-кеша", frequencyArray[0].fileName)
	return nil
}
