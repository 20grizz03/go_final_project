package repeatTask

import (
	"errors"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату на основе правила повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// откуда начинается отсчет
	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", errors.New("некорректный формат даты")
	}

	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	rules := strings.Split(repeat, " ")
	switch rules[0] {
	// правило для дней
	case "d":

		return DailyRepeat(now, taskDate, rules)
	// правило для года
	case "y":
		return YearRepeat(now, taskDate, rules)
	// правило для недели
	case "w":
		return WeekRepeat(now, taskDate, rules)
	// правило для месяца
	case "m":
		return MonthRepeat(now, taskDate, rules)
	default:
		return "", errors.New("правило повторения не поддерживается")
	}
}
