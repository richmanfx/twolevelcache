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
func dataPlotting(data []int64, ramCacheSize, driveCacheSize, dataNumber, requestAmount int) {

	log.Debugf("Число точек: %d", len(data))
	log.Debugf("Данные: %v", data)

	plt, err := plot.New()
	if err != nil {
		panic(err)
	}

	plt.Title.Text = fmt.Sprintf(
		"Задержка получения данных\n"+
			"Размер RAM-кеша: %d, размер DRIVE-кеша: %d, количество данных: %d, количество запросов: %d",
		ramCacheSize, driveCacheSize, dataNumber, requestAmount)
	plt.X.Label.Text = "Порядковый номер случайного запроса"
	plt.Y.Label.Text = "Задержка, нс"

	// Сетка
	plt.Add(plotter.NewGrid())

	// Добавить на график линию с точками из координат
	err = plotutil.AddLinePoints(plt, "", getPoints(data))
	if err != nil {
		panic(err)
	}

	// Записать результаты в графический файл
	if err := plt.Save(18*vg.Inch, 10*vg.Inch, graphResultFileName); err != nil {
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
