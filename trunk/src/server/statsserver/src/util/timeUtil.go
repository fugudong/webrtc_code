package util


const (
	// 定义每分钟的秒数
	SecondsPerMinute = 60
	// 定义每小时的秒数
	SecondsPerHour = SecondsPerMinute * 60
	// 定义每天的秒数
	SecondsPerDay = SecondsPerHour * 24
)

func ResolveTime(seconds int64) (hour int, minute int, seccond int) {
	hour = int(seconds / SecondsPerHour)
	minute = int(seconds / SecondsPerMinute)
	seccond = int(seconds % SecondsPerMinute)
	return
}
