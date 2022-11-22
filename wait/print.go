package wait

import "fmt"

var CurrentOutputLevel = INFO

type OutputLevel int

const (
	SILENT OutputLevel = iota
	INFO
)

func Println(a ...any) {
	if CurrentOutputLevel == INFO {
		fmt.Println(a...)
	}
}
