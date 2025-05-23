package main

import "fmt"

func main() {
	var message string
	fmt.Print("Enter decoded message: ")
	fmt.Scanln(&message)
	fmt.Println("Decode ways:", decode(message))
}

func decode(message string) int {
	if message == "" {
		return 0
	}
	if message[0] == '0' {
		return 0
	}
	current := 1
	prev := 1
	prev_prev := 1
	var char byte
	var prev_char byte

	for i := 1; i < len(message); i++ {
		char = message[i]
		prev_char = message[i-1]
		if char == '0' {
			if prev_char == '0' || prev_char > '2' {
				return 0
			}
			current = prev_prev
		} else if prev_char == '0' {
			current = prev
		} else if string(prev_char)+string(char) > "26" {
			current = prev
		} else {
			current = prev + prev_prev
		}
		prev_prev = prev
		prev = current
	}
	return current
}
