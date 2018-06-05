package show

import (
	"fmt"
	"sql-client/pkg/transfer"

	"github.com/mssola/colors"
)

var color *colors.Color = nil

func init() {
	color = colors.Default()
}

type ShowInfo struct {
	Info  string
	Color colors.Colors
}

func New(info string, red colors.Colors) *ShowInfo {
	return &ShowInfo{
		Info:  info,
		Color: red,
	}
}

func Println(info *ShowInfo) {
	if nil == info || "" == info.Info {
		return
	}

	color.Change(info.Color, false)
	fmt.Println(color.Get(info.Info))
}

func TitlePrintln(info *ShowInfo) {
	if nil == info || "" == info.Info {
		return
	}

	fmt.Print(">>> ")
	color.Change(info.Color, false)
	fmt.Println(color.Get(info.Info))
}

func PrintListln(infos []*ShowInfo) {
	var tmp string = ""
	for _, info := range infos {
		if nil == info || "" == info.Info {
			continue
		}
		color.Change(info.Color, false)
		tmp += color.Get(info.Info)
	}
	if "" != tmp {
		fmt.Println(tmp)
	}
	return
}

func TitlePrintListln(infos []*ShowInfo) {
	fmt.Print(">>> ")
	var tmp string = ""
	for _, info := range infos {
		if nil == info || "" == info.Info {
			continue
		}
		color.Change(info.Color, false)
		tmp += color.Get(info.Info)
	}
	if "" != tmp {
		fmt.Println(tmp)
	}
	return
}

func Title(num int64) {
	Println(New(fmt.Sprintf("rows affected num: %d.", num), colors.Red))
}

func Header(fields []string) {
	var infos []*ShowInfo = nil
	infos = append(infos, New("table fields: ", colors.Red))
	for i, field := range fields {
		infos = append(infos, New(fmt.Sprintf("  %s  ", field), colors.Red))
		if i != len(fields)-1 {
			infos = append(infos, New("|", colors.Blue))
		}
	}
	if len(infos) > 0 {
		PrintListln(infos)
	}
	return
}

func Body(num int, values []interface{}, types []string) {
	var infos []*ShowInfo = nil
	infos = append(infos, New(fmt.Sprintf("row num:%d ", num), colors.Red))
	for i, value := range values {
		if nil == value {
			continue
		}
		if nil != types {
			infos = append(infos, New(fmt.Sprintf(" %s ", transfer.ToString(types[i], value)), colors.Green))
		} else {
			infos = append(infos, New(fmt.Sprintf(" %v ", value), colors.Green))
		}
		if i != len(values)-1 {
			infos = append(infos, New("|", colors.Blue))
		}
	}
	if len(infos) > 0 {
		TitlePrintListln(infos)
	}

	return
}
