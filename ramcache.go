/* Реализация кеша в оперативной памяти */

package main

import (
	"errors"
	"sync"
)

/* Контейнер */
type MemoryCache struct {
	sync.RWMutex                           // Защитный мьютекс
	elements     map[string]*MemoryElement // Кешируемые элементы

}

/* Кешируемый элемент */
type MemoryElement struct {
	value     interface{} // Кешируемое значение
	frequency int         // Частота использорвания элемента
}

/* Создать новый кеш в памяти */
func CreateMemoryCache() Cache {
	return &MemoryCache{elements: make(map[string]*MemoryElement)}
}

/***************************************/
/* Реализация методов интерфейса Cache */

/* Put */
func (mc *MemoryCache) Put(key string, value interface{}) error {

	// Блокировать на время записи
	mc.Lock()
	defer mc.Unlock()

	mc.elements[key] = &MemoryElement{
		value:     value,
		frequency: 1, // Помещаем в кеш - значит используется в первый раз
	}
	return nil
}

/* Get */
func (mc *MemoryCache) Get(key string) interface{} {

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	// Если нет значения в кеше, то вернуть nil
	el, ok := mc.elements[key]
	if ok {
		el.frequency++ // Частота использования
		return el.value
	} else {
		return nil
	}
}

/* Del */
func (mc *MemoryCache) Del(key string) error {

	// Блокировать на время записи
	mc.Lock()
	defer mc.Unlock()

	_, ok := mc.elements[key]
	if !ok {
		return errors.New("удаление значения: ключ не существует")
	}

	delete(mc.elements, key)
	_, ok = mc.elements[key]
	if ok {
		return errors.New("ошибка при удалении значения")
	}

	return nil
}

/* IsExist */
func (mc *MemoryCache) IsExist(key string) bool {

	// Блокировать запись на время чтения
	mc.RLock()
	defer mc.RUnlock()

	_, ok := mc.elements[key]
	if ok {
		return true
	} else {
		return false
	}
}
