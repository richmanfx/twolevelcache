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
	value     interface{} // Кешируемое значение
	frequency int         // Частота использования элемента
	// (может использовать TreeMap???
	// import "github.com/golang-collections/tree/treemap"
}

/* Создать новый кеш в памяти без ограничения размера */
func CreateMemoryCache() Cache {
	return &MemoryCache{elements: make(map[string]*MemoryElement), maxSize: -1}
}

/* Создать новый кеш в памяти заданного размера */
func CreateSpecifySizeMemoryCache(size int) Cache {
	return &MemoryCache{elements: make(map[string]*MemoryElement), maxSize: size}
}

/***************************************/
/* Реализация методов интерфейса Cache */

/* Put */
func (mc *MemoryCache) Put(key string, value interface{}) error {

	// Блокировать на время записи
	mc.Lock()
	defer mc.Unlock()

	// Проверить не заполнен ли кеш полностью
	if mc.maxSize != -1 { // "-1" - нет ограничения в размере кеша
		log.Infof("Количество элементов в кеше: %d", mc.Size())
		log.Infof("Максимальный размер кеша: %d", mc.maxSize)
		if mc.Size() >= mc.maxSize {

			// Передвигать всё на диск	// TODO

			log.Infoln("Кеш полностью заполнет, ничего не добавляем!")
			return errors.New("кеш полностью заполнет, ничего не добавляем")
		}
	}

	mc.elements[key] = &MemoryElement{
		value:     value,
		frequency: 1, // Помещаем в кеш - значит используется в первый раз
	}

	return nil
}

/* Get */
func (mc *MemoryCache) Get(key string) interface{} {

	var result interface{}

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	// Если нет значения в кеше, то вернуть nil
	el, ok := mc.elements[key]
	if ok {
		el.frequency++ // Частота использования
		result = el.value
	} else {
		result = nil
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
	} else {
		result = false
	}

	return result
}

/* Size */
func (mc *MemoryCache) Size() int {
	result := len(mc.elements)
	log.Debugf("Количество элементов в кеше: %d", result)
	return result
}
