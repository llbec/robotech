package contract

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type stateMutability = int

const (
	View       stateMutability = 0
	Nonpayable stateMutability = 1
)

type ABI struct {
	Name            string `json:"name"`
	StateMutability string `json:"stateMutability"`
	Type            string `json:"type"`
	//Inputs
}

func ReadInterfaces(_path string) error {
	f, err := ioutil.ReadFile(_path)
	if err != nil {
		return err
	}
	contractABI, err = getFunction(string(f))
	return err
}

func getFunction(_context string) (abiStr string, err error) {
	abiStr = "{\n"
	//type function or event in the head
	regHead := regexp.MustCompile(`^[a-z]+`)
	if regHead == nil {
		return "", fmt.Errorf("type MustCompile error")
	}
	headIndex := regHead.FindStringIndex(_context)
	abiType := _context[headIndex[0]:headIndex[1]]

	//verify and stateMutability
	if abiType == "function" {
		if !strings.Contains(_context, "public") && !strings.Contains(_context, "external") {
			return "", fmt.Errorf("not a public or external function")
		}
		if strings.Contains(_context, "view") {
			abiStr += "\t\"stateMutability\": \"view\",\n"
		} else {
			abiStr += "\t\"stateMutability\": \"nonpayable\",\n"
		}
	}

	//inputs and outputs
	regPuts := regexp.MustCompile(`\(([^)]*)\)`)
	if regPuts == nil {
		return "", fmt.Errorf("puts MustCompile error")
	}
	puts := regPuts.FindAllString(_context, -1)
	if len(puts) != 1 && len(puts) != 2 {
		return "", fmt.Errorf("invalid puts: %v", puts)
	}
	//inputs
	if len(puts[0]) == 2 {
		abiStr += "\t\"inputs\": [],\n"
	} else {
		abiStr += "\t\"inputs\": [\n"
		inputs := strings.Split(puts[0][1:len(puts[0])-1], ",")
		for i, input := range inputs {
			abiStr += "\t\t{\n"
			rs := strings.Split(input, " ")
			num := len(rs)
			if num == 2 {
				internalType, t := getArgType(rs[0])
				abiStr += fmt.Sprintf("\t\t\t\"internalType\": \"%s\",\n", internalType)
				abiStr += fmt.Sprintf("\t\t\t\"name\": \"%s\",\n", rs[1])
				abiStr += fmt.Sprintf("\t\t\t\"type\": \"%s\"\n", t)
			} else if num == 3 {
				panic("struct wait for handle")
			} else {
				return "", fmt.Errorf("invalid input: %v", input)
			}
			if i == len(inputs)-1 {
				abiStr += "\t\t}\n"
			} else {
				abiStr += "\t\t},\n"
			}
		}
		abiStr += "\t],\n"
	}
	//outputs
	if len(puts) != 2 {
		abiStr += "\t\"outputs\": [],\n"
	} else {
		abiStr += "\t\"outputs\": [\n"
		outputs := strings.Split(puts[1][1:len(puts[1])-1], ",")
		for i, output := range outputs {
			abiStr += "\t\t{\n"
			if strings.Contains(output, " ") {
				panic(fmt.Errorf("struct wait for handle: %v", output))
			} else {
				internalType, t := getArgType(output)
				abiStr += fmt.Sprintf("\t\t\t\"internalType\": \"%s\",\n", internalType)
				abiStr += "\t\t\t\"name\": \"\",\n"
				abiStr += fmt.Sprintf("\t\t\t\"type\": \"%s\"\n", t)
			}
			//
			if i == len(outputs)-1 {
				abiStr += "\t\t}\n"
			} else {
				abiStr += "\t\t},\n"
			}
		}
		//
		abiStr += "\t],\n"
	}

	//name
	regName := regexp.MustCompile(` \S+\(`)
	if regName == nil {
		return "", fmt.Errorf("name MustCompile error")
	}
	nameIndex := regName.FindStringIndex(_context)
	abiStr += fmt.Sprintf("\t\"name\": \"%s\",\n", _context[nameIndex[0]+1:nameIndex[1]-1])

	abiStr += fmt.Sprintf("\t\"type\": \"%s\"\n", abiType)
	abiStr += "}"
	return
}

func getArgType(_type string) (string, string) {
	regInt := regexp.MustCompile(`^int[0-9]*`)
	regUint := regexp.MustCompile(`^uint[0-9]*`)
	regByte := regexp.MustCompile(`^byte[0-9]*`)
	keys := make(map[string]bool)
	keys["bool"] = true
	keys["address"] = true
	keys["string"] = true
	if regInt == nil || regUint == nil || regByte == nil {
		goto NORESULT
	}
	if regInt.MatchString(_type) {
		return _type, _type
	} else if regUint.MatchString(_type) {
		return _type, _type
	} else if regByte.MatchString(_type) {
		return _type, _type
	} else if keys[_type] {
		return _type, _type
	} else {
		return fmt.Sprintf("contract %s", _type), "address"
	}
NORESULT:
	return "", ""
}
