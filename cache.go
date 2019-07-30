package main

/* Базовый интерфейс кеша */
type Cache interface {

	// Поместить кешируемыое значение с ключом в кеш
	Put(key string, val interface{}) error

	// Получить значение из кеша по ключу
	Get(key string) *MemoryElement

	// Удалить значение из кеша по ключу
	Del(key string) error

	// Проверить существование значения в кеше
	IsExist(key string) bool

	// Размер заполнения кеша
	Size() int
}
