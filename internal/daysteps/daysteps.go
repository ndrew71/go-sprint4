package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	s := strings.Split(data, ",")
	if len(s) != 2 {
		return 0, 0, fmt.Errorf("нужно два значения")
	}

	n, err := strconv.Atoi(s[0])
	if err != nil {
		return 0, 0, err
	}
	if n <= 0 {
		return 0, 0, fmt.Errorf("шагов должно быть больше нуля")
	}

	t,err := time.ParseDuration(s[1])
	if err != nil {
		return 0, 0, err
	}
	if t <= 0 {
		return 0,0, fmt.Errorf("время должно быть больше нуля")
	}
	return n, t, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, t, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if steps <= 0 {
		return ""
	}

	d := float64(steps)*stepLength / mInKm
	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, t)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps,d,cal)
}
