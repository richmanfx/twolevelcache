package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

/* Вывести график в файл */
func dataPlotting(data []int64, cacheSize int, dataAmount int, requestAmount int) {

	log.Debugf("Число точек: %d", len(data))
	log.Debugf("Данные: %v", data)

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = fmt.Sprintf(
		"Задержка получения данных.\nРазмер кеша: %d, количество данных: %d, количество запросов: %d",
		cacheSize, dataAmount, requestAmount)
	plt.X.Label.Text = "Порядковый номер случайного запроса"
	plt.Y.Label.Text = "Задержка"

	plt.Add(plotter.NewGrid()) // Сетка

	err = plotutil.AddLinePoints(plt, "", getPoints(data))
	if err != nil {
		panic(err)
	}

	// Записать график в PNG файл
	if err := plt.Save(18*vg.Inch, 10*vg.Inch, "request_delay.jpg"); err != nil { // TODO: имя файла и размер - в конфиг
		panic(err)
	}
}

/* Вернуть координаты точек графика, вычисленных на основе данных */
func getPoints(data []int64) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = float64(data[i])
	}
	return pts
}
