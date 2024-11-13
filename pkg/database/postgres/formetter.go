package postgres

import "time"

type ConditionSign string

const (
	Gt  ConditionSign = ">"
	Lt  ConditionSign = "<"
	Gte ConditionSign = ">="
	Lte ConditionSign = "<="
)

type TimeCondition map[ConditionSign]time.Time

func FormatTimeCondition(sign ConditionSign, date time.Time) TimeCondition {
	return TimeCondition{
		sign: date,
	}
}
