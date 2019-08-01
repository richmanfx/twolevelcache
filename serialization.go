package main

import (
	"encoding/gob"
	"errors"
	log "github.com/Sirupsen/logrus"
	"os"
)

/* Сериализация в двоичном формате Go */
func gobEncode(fileName string, data interface{}) error {

	gob.Register(SimpleStructure{}) // Регистрация типа

	fullPath := cacheDir + "/" + fileName
	file, err := os.Create(fullPath)
	if err != nil {
		log.Errorf("Ошибка создания файла кеширования '%s': %s", fullPath, err)
		return errors.New("ошибка создания файла кеширования")
	}

	encoder := gob.NewEncoder(file)

	err = encoder.Encode(data)
	if err != nil {
		log.Errorf("Ошибка кодирования: %s", err)
		return errors.New("ошибка 'mob' кодирования")
	}

	err = file.Close()
	if err != nil {
		log.Errorf("Ошибка закрытия файла кеширования '%s': %s", fullPath, err)
		return errors.New("ошибка закрытия файла кеширования")
	}

	return nil
}

/* Десериализация из двоичного формата Go */
func gobDecode(fileName string) (data *MemoryElement, err error) {

	fullPath := cacheDir + "/" + fileName
	file, err := os.Open(fullPath)
	if err != nil {
		log.Errorf("Ошибка открытия файла кеширования '%s': %s", fullPath, err)
		return nil, errors.New("ошибка открытия файла кеширования")
	}

	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&data)
	if err != nil {
		log.Errorf("Ошибка декодирования: %s", err)
		return nil, errors.New("ошибка 'mob' декодирования")
	}

	err = file.Close()
	if err != nil {
		log.Errorf("Ошибка закрытия файла кеширования '%s': %s", fullPath, err)
		return nil, errors.New("ошибка закрытия файла кеширования")
	}

	return data, nil
}
