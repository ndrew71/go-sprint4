package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 3 {
		return 0,"",0,fmt.Errorf("нужно три значения")
	}

	n, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, "", 0, err
	}
	if n <= 0 {
		return 0, "", 0, fmt.Errorf("шагов должно быть больше нуля")
	}
	t, err := time.ParseDuration(s[2])
	if err != nil {
		return 0, "", 0, err
	}
	if t <= 0 {
		return 0, "", 0, fmt.Errorf("время должно быть больше нуля")
	}

	return n,s[1],t,nil
}

func distance(steps int, height float64) float64 {
	step := height * stepLengthCoefficient
	return float64(steps) * step / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	d := distance(steps,height)
	return d / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, kind, t, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "",err
	}

	var cal float64
	switch kind {
	case "Бег":
		cal, err = RunningSpentCalories(steps, weight, height, t)
	case "Ходьба":
		cal, err = WalkingSpentCalories(steps, weight, height, t)
	default:
		return "", fmt.Errorf("неизвестный тип тренировки")
	}
	if err != nil {
		return "", err
	}

	d := distance(steps,height)
	v := meanSpeed(steps, height, t)
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", kind, t.Hours(), d, v, cal), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("неверные данные тренировки")
	}
	v := meanSpeed(steps,height,duration)
	cal := weight * v * duration.Minutes() / minInH
	return cal,nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("неверные данные прогулки")
	}

	v := meanSpeed(steps, height, duration)
	cal := weight*v*duration.Minutes() / minInH
	return cal * walkingCaloriesCoefficient, nil
}
