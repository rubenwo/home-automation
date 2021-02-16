package rlang

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Script func()

type keyword string

const (
	turnOn     keyword = "turn_on"
	turnOff    keyword = "turn_off"
	blockStart keyword = "<>"
	blockEnd   keyword = "</>"
)

type instruction struct {
	keyword keyword
	args    []interface{}
}

//Lang-example

// <>
// turn_on>
// </>
//

func tokenize(body string) ([]instruction, error) {
	scanner := bufio.NewScanner(bytes.NewBufferString(body))
	var instructions []instruction
	var tempInstruction instruction

	for scanner.Scan() {
		line := scanner.Text()
		tempInstruction = instruction{}
		split := strings.Split(line, ">")
		switch split[0] {
		case "turn_on":
			tempInstruction.keyword = turnOn
			tempInstruction.args = append(tempInstruction.args, split[1])
		case "turn_off":
			tempInstruction.keyword = turnOff
			tempInstruction.args = append(tempInstruction.args, split[1])
		}
		instructions = append(instructions, tempInstruction)
	}

	return instructions, nil
}

func Parse(body string) (Script, error) {
	fmt.Println("parsing")
	instructions, err := tokenize(body)
	if err != nil {
		return nil, err
	}
	fmt.Println(instructions, err)
	return func() {
		fmt.Println("Executing function")
		for i, instruction := range instructions {
			fmt.Println("running instruction at index:", i)
			switch instruction.keyword {
			case turnOn:
				client := &http.Client{}
				req, err := http.NewRequest("GET", instruction.args[0].(string), nil)
				if err != nil {
					log.Println(err)
					return
				}
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
					return
				}
				defer resp.Body.Close()
			}

		}
	}, nil
}
