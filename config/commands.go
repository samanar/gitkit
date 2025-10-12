package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func Ask(question, def string) string {
	fmt.Printf("%s [%s]: ", question, def)
	answer := ReadLine()
	if answer == "" {
		return def
	}
	return answer
}
