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

	// Проверить не заполнен ли RAM-кеш полностью
	if mc.maxSize != -1 { // "-1" - нет ограничения в размере кеша
		log.Debugf("Количество элементов в RAM-кеше: %d", mc.Size())
		log.Debugf("Максимальный размер RAM-кеша: %d", mc.maxSize)

		if mc.Size() >= mc.maxSize {
			log.Infoln("RAM-кеш полностью заполнен - рекешируем")

			// Рекеширование
			err := reCaching(mc)
			if err != nil {
				log.Infof("Ошибка рекеширования: %s", err)
			}

		}
	}

	// Поместить в RAM-кеш
	mc.elements[key] = value

	return nil
}

/* Get */
func (mc *MemoryCache) Get(key string) *MemoryElement {

	var result *MemoryElement

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	// Если нет значения в RAM-кеше, то вернуть "nil".
	element, ok := mc.elements[key]
	if ok {
		element.Frequency++ // Частота использования
		result = element
	} else {
		result = nil // Нет данных в RAM-кеше
	}

	return result
}

/* Del */
func (mc *MemoryCache) Del(key string) error {

	var result error

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
	log.Debugf("Количество элементов в RAM-кеше: %d", result)
	return result
}
