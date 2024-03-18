package exception

import "fmt"

func ReportError(err error, message string) {
	if err != nil {
		fmt.Printf("%s : %v", message, err)
	}
}
