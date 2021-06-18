package cmd

import "fmt"

const banner = `
                               _     _
 _ __   ___   ___ __ _ ___ ___(_)___| |_
| '_ \ / _ \ / __/ _' / __/ __| / __| __|
| |_) | (_) | (_| (_| \__ \__ \ \__ \ |_
| .__/ \___/ \___\__,_|___/___/_|___/\__|
|_|
`

func PrintBanner() {
	fmt.Printf("%s\n", banner)
	fmt.Printf("\t\tv1.0.0\n\n")
	fmt.Printf("\t\thttps://pocassist.jweny.top/\n\n")
}

