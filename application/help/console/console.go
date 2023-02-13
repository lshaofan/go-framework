package console

const (
	darkgreen    = "\033[36;41m"
	whiteBlack   = "\033[30;47m"
	yellowBlue   = "\033[34;43m"
	redBlue      = "\033[34;41m"
	blueYellow   = "\033[33;44m"
	magentaGreen = "\033[32;45m"
	cyanYellow   = "\033[97;46m"
	greenWhite   = "\u001B[97;42m"
	reset        = "\033[0m"

	green   = "\033[32m"
	white   = "\033[37m"
	yellow  = "\033[33m"
	red     = "\033[31m"
	blue    = "\033[34m"
	magenta = "\033[35m"
	cyan    = "\033[36m"
)

// 格式化输出

func DarkGreen(s string) {
	println(darkgreen, s, reset)
}

func WhiteBlack(s string) {
	println(whiteBlack, s, reset)
}

func YellowBlue(s string) {
	println(yellowBlue, s, reset)
}

func RedBlue(s string) {
	println(redBlue, s, reset)
}

func BlueYellow(s string) {
	println(blueYellow, s, reset)
}

func MagentaGreen(s string) {
	println(magentaGreen, s, reset)
}

func CyanBg(s string) {
	println(cyanYellow, s, reset)
}

func GreenWhite(s string) {
	println(greenWhite, s, reset)
}
func Green(s string) {
	println(green, s, reset)
}

func White(s string) {
	println(white, s, reset)
}

func Yellow(s string) {
	println(yellow, s, reset)
}

func Red(s string) {
	println(red, s, reset)
}

func Blue(s string) {
	println(blue, s, reset)
}

func Magenta(s string) {
	println(magenta, s, reset)
}

func Cyan(s string) {
	println(cyan, s, reset)
}

func Error(s string) {
	println(red, s, reset)
}

func Success(s string) {
	println(green, s, reset)
}

func Info(s string) {
	println(blue, s, reset)
}

// 用法示例
// console.DarkGreen("DarkGreen")
