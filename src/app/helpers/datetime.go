package helpers

import "time"

const ISO8601 = "2006-01-02T15:04:05-0700"

func GetDateForDB(t time.Time) time.Time {
	return t.UTC()
}

func GetDateForJSON(t time.Time) string {
	return t.UTC().Format(ISO8601)
}
