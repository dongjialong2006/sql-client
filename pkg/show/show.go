package show

import (
	"fmt"

	"github.com/mssola/colors"
)

func println(info string, color *colors.Color) {
	if "" != info {
		fmt.Println(color.Get(info))
	}
}

func Header(fields []string, color *colors.Color) {
	fmt.Print(">>> ")
	color.Change(colors.Red, false)
	var tmp string = color.Get("table fields: ")
	for i, field := range fields {
		color.Change(colors.Red, false)
		tmp += color.Get(fmt.Sprintf("  %s  ", field))
		if i != len(fields)-1 {
			color.Change(colors.Blue, false)
			tmp += color.Get("|")
		}
	}
	println(tmp, color)
	return
}

func Body(num int, values []interface{}, color *colors.Color) {
	color.Change(colors.Red, false)
	var tmp string = color.Get(fmt.Sprintf("    row num:%d ", num))
	fmt.Print(tmp)

	tmp = ""
	for i, value := range values {
		color.Change(colors.Green, false)
		tmp += color.Get(fmt.Sprintf(" %v ", value))
		if i != len(values)-1 {
			color.Change(colors.Blue, false)
			tmp += color.Get("|")
		}
	}
	println(tmp, color)
	return
}
