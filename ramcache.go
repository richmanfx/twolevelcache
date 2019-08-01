/* Реализация кеша в оперативной памяти */

package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"sync"
)

/* Контейнер */
type MemoryCache struct {
	sync.RWMutex                           // Защитный мьютекс для потокобезопасности map
	elements     map[string]*MemoryElement // Кешируемые элементы
	maxSize      int                       // Максимальный размер кеша
}

/* Кешируемый элемент */
type MemoryElement struct {
	Value     interface{} // Кешируемое значение
	Frequency int         // Частота использования элемента
}

/* Создать новый кеш в памяти заданного размера */
func CreateSpecifySizeRamCache(size int) *MemoryCache {
	return &MemoryCache{elements: make(map[string]*MemoryElement), maxSize: size}
}

/* Реализация методов интерфейса Cache */

/* Put */
func (mc *MemoryCache) Put(key string, value *MemoryElement) error {

	// Блокировать на время записи
	mc.Lock()
	defer mc.Unlock()

	// Проверить не заполнен ли кеш полностью
	if mc.maxSize != -1 { // "-1" - нет ограничения в размере кеша
		log.Debugf("Количество элементов в кеше: %d", mc.Size())
		log.Debugf("Максимальный размер кеша: %d", mc.maxSize)

		if mc.Size() >= mc.maxSize {
			log.Infoln("RAM-кеш полностью заполнен - рекешируем")

			// Рекеширование
			err := reCaching()
			if err != nil {
				log.Infof("Ошибка рекеширования: %s", err)
			}

		}
	}

	// Поместить в ram-кеш
	mc.elements[key] = value

	return nil
}

/* Get */
func (mc *MemoryCache) Get(key string) *MemoryElement {

	var result *MemoryElement

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	// Если нет значения в ram-кеше, то проверить в drive-кеше. Если и там нет, то вернуть "nil".
	element, ok := mc.elements[key]
	if ok {
		element.Frequency++ // Частота использования
		result = element
	} else {
		// Проверить в drive-кеше
		// TODO: Нужно ли здесь проверять drive-кеш???

		result = nil // Нигде нет
	}

	return result
}

/* Del */
func (mc *MemoryCache) Del(key string) error {

	var result error

	// Блокировать на время записи
	mc.Lock()
	defer mc.Unlock()

	_, ok := mc.elements[key]
	if !ok {
		result = errors.New("удаление значения: ключ не существует")
	} else {
		delete(mc.elements, key)
		_, ok = mc.elements[key]
		if ok {
			result = errors.New("ошибка при удалении значения")
		} else {
			result = nil
		}
	}
	return result
}

/* IsExist */
func (mc *MemoryCache) IsExist(key string) bool {

	var result bool

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	_, ok := mc.elements[key]
	if ok {
		result = true
	}

	return result
}

/* Size */
func (mc *MemoryCache) Size() int {
	result := len(mc.elements)
	log.Debugf("Количество элементов в ram-кеше: %d", result)
	return result
}
