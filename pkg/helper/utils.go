package helper

import (
	"fmt"
	"github.com/fatih/color"
)

// ASCII art logo
const ConsoleStartStr = `                _        _     _         _ 
               | |      | |   (_)       | |
   ___ ___   __| | ___  | |__  _ _ __ __| |
  / __/ _ \ / _` + "`" + ` |/ _ \ | '_ \| | '__/ _` + "`" + ` |
 | (_| (_) | (_| |  __/ | |_) | | | | (_| |
  \___\___/ \__,_|\___| |_.__/|_|_|  \__,_|
                                           
                                           `

// 打印logo
func PrintLogo() {
	// 打印成绿色的
	fmt.Println(Green(ConsoleStartStr))
}

// Green 绿色
func Green(s string) string {
	return color.GreenString(s)
}
