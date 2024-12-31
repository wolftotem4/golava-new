package question

import (
	"fmt"
)

func AskOverwrite(files []string) bool {
	for _, file := range files {
		fmt.Println(file)
	}

	fmt.Print("Above files will be overwritten. Do you want to continue? (y/n): ")
	var answer string
	fmt.Scanln(&answer)
	return answer == "y"
}
