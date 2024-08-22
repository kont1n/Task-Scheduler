package internal

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	conditions := strings.Split(repeat, " ")
	switch conditions[0] {
	case "d":
		{
			if len(repeat) <= 1 {
				return "", errors.New("не указан интервал в днях")
			}

			i, err := strconv.Atoi(conditions[1])
			if err != nil {
				return "", err
			}

			if i > 400 {
				return "", errors.New("превышение максимального допустимого числа дней")
			}

			newDate, err := DayTransfer(now, i, date)
			if err != nil {
				return "", err
			}
			return newDate, nil
		}
	case "y":
		{
			newDate, err := YearTransfer(now, date)
			if err != nil {
				return "", err
			}
			return newDate, nil
		}
	case "w":
		{
			if len(repeat) <= 1 {
				return "", errors.New("недопустимое значение")
			}
			weekDays := strings.Split(conditions[1], ",")
			var d = []int{}
			for _, i := range weekDays {
				j, err := strconv.Atoi(i)
				if err != nil {
					return "", err
				}
				if j > 31 {
					return "", errors.New("недопустимый день месяца")
				}
				d = append(d, j)
			}
			newDate, err := WeekDistribution(now, d)
			if err != nil {
				return "", err
			}
			return newDate, nil
		}
	case "m":
		{
			if len(repeat) <= 1 {
				return "", errors.New("недопустимое значение")
			}
			monthDays := strings.Split(conditions[1], ",")
			var d = []int{}
			for _, i := range monthDays {
				j, err := strconv.Atoi(i)
				if err != nil {
					return "", err
				}
				if j > 31 {
					return "", errors.New("недопустимый день месяца")
				}
				d = append(d, j)
			}

			if conditions[2] != "" {

				months := strings.Split(conditions[2], ",")
				var m = []int{}
				for _, i := range months {
					j, err := strconv.Atoi(i)
					if err != nil {
						return "", err
					}
					if j > 12 {
						return "", errors.New("недопустимый месяц")
					}
					m = append(m, j)
				}

				newDate, err := MonthDistribution(now, d, m)
				if err != nil {
					return "", err
				}
				return newDate, nil
			} else {
				newDate, err := MonthDistribution(now, d, nil)
				if err != nil {
					return "", err
				}
				return newDate, nil
			}
		}
	default:
		{
			return "", errors.New("правило повторения указано в неправильном формате")
		}
	}
}

func DayTransfer(now time.Time, step int, date string) (string, error) {
	var newDay time.Time
	var err error

	newDay, err = time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	for {
		newDay = newDay.AddDate(0, 0, step)
		if newDay.After(now) {
			break
		}
	}
	return newDay.Format("20060102"), nil
}

func YearTransfer(now time.Time, date string) (string, error) {
	var newYear time.Time
	var err error

	newYear, err = time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	for {
		newYear = newYear.AddDate(1, 0, 0)
		if newYear.After(now) {
			break
		}
	}
	return newYear.Format("20060102"), nil
}

func WeekDistribution(now time.Time, days []int) (string, error) {
	var newDay time.Time
	var dates []time.Time

	for _, day := range days {
		d, err := WeekDays(day)
		if err != nil {
			return "", err
		}

		number := int(now.Weekday())

		if d >= number {
			newDay = now.AddDate(0, 0, d-number)
		} else {
			newDay = now.AddDate(0, 0, 7)
		}

		dates = append(dates, newDay)
	}

	// TBD: date sort

	return dates[0].Format("20060102"), nil
}

func MonthDistribution(now time.Time, days []int, months []int) (string, error) {
	var newDay time.Time
	var dates []time.Time

	y, m, d := now.Date()

	if len(months) == 0 {
		for _, day := range days {
			if day == -1 {
				newDay = time.Date(y, m+1, -1, 0, 0, 0, 0, time.Local)
			} else if day == -2 {
				newDay = time.Date(y, m+1, -2, 0, 0, 0, 0, time.Local)
			} else {
				if day > d {
					newDay = time.Date(y, m, day, 0, 0, 0, 0, time.Local)
				} else {
					newDay = time.Date(y, m+1, day, 0, 0, 0, 0, time.Local)
				}
			}
			dates = append(dates, newDay)
		}
	} else {
		for _, month := range months {
			for _, day := range days {
				if day == -1 {
					newDay = time.Date(y, time.Month(month+1), -1, 0, 0, 0, 0, time.Local)
				} else if day == -2 {
					newDay = time.Date(y, time.Month(month+1), -2, 0, 0, 0, 0, time.Local)
				} else {
					if day > d {
						newDay = time.Date(y, time.Month(month), day, 0, 0, 0, 0, time.Local)
					} else {
						newDay = time.Date(y, time.Month(month+1), day, 0, 0, 0, 0, time.Local)
					}
				}
				dates = append(dates, newDay)
			}
		}
	}

	// TBD: date sort

	return dates[0].Format("20060102"), nil

}

func WeekDays(d int) (int, error) {
	switch d {
	case 1:
		return int(time.Monday), nil
	case 2:
		return int(time.Tuesday), nil
	case 3:
		return int(time.Wednesday), nil
	case 4:
		return int(time.Tuesday), nil
	case 5:
		return int(time.Friday), nil
	case 6:
		return int(time.Saturday), nil
	case 7:
		return int(time.Sunday), nil
	default:
		return 0, errors.New("weekdays error")
	}
}
