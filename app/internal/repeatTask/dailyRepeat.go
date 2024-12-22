package repeatTask

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func DailyRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	// проверяем правильность структуры правила
	if len(rules) != 2 {
		return "", errors.New("неверный формат правила повторения для дней")
	}

	// преобразуем количество дней в число
	days, err := strconv.Atoi(rules[1])
	if err != nil || days <= 0 || days > 400 {
		return "", errors.New("некорректный интервал для правила повторения дней")
	}

	if taskDate.After(now) {
		taskDate = taskDate.AddDate(0, 0, days)
	} else {
		// Добавляем дни до тех пор, пока `taskDate` не станет больше `now`
		for !taskDate.After(now) {
			taskDate = taskDate.AddDate(0, 0, days)
		}
	}
	// Возвращаем следующую дату в формате YYYYMMDD
	return taskDate.Format("20060102"), nil
}

func parseWeekdays(days string) ([]time.Weekday, error) {
	dayStrings := strings.Split(days, ",")
	var weekdays []time.Weekday

	for _, dayStr := range dayStrings {
		dayInt, err := strconv.Atoi(dayStr)
		if err != nil || dayInt < 1 || dayInt > 7 {
			return nil, errors.New("некорректный день недели в правиле повторения")
		}
		weekday := time.Weekday((dayInt % 7))
		weekdays = append(weekdays, weekday)
	}
	return weekdays, nil
}

// обработчик дня
func WeekRepeat(now time.Time, taskDate time.Time, rules []string) (string, error) {
	// проверяем правильность структуры правила
	if len(rules) != 2 {
		return "", errors.New("неверный формат правила повторения для недели")
	}
	weekdays, err := parseWeekdays(rules[1])
	if err != nil {
		return "", err
	}
	for {
		// проверяем, совпадает ли текущий день недели с одним из указанных
		for _, weekday := range weekdays {
			if taskDate.Weekday() == weekday {
				if taskDate.After(now) {
					return taskDate.Format("20060102"), nil
				}
			}
		}
		// Если не совпало, переносим дату на следующий день
		taskDate = taskDate.AddDate(0, 0, 1)
	}
}
