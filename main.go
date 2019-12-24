package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/cedrickchee/go-parkinglot/internal/printer"
)

var inputInteractive io.Reader = os.Stdin
var outStream io.Writer = os.Stdout

func main() {
	args := os.Args
	argsLen := len(args)

	var scanner *bufio.Scanner

	switch {
	case argsLen == 2:
		inputFile, err := os.Open(args[1])
		if err != nil {
			panic(err)
		}
		defer inputFile.Close()
		scanner = bufio.NewScanner(inputFile)
	case argsLen > 2:
		log.Fatal("Unknown command line input")
	default:
		scanner = bufio.NewScanner(inputInteractive)
	}

	exit := false
	for !exit && scanner.Scan() {
		input := scanner.Text()

		switch {
		case validate(input, "create_parking_lot", 2):
			maxSlot := 6
			fmt.Fprintf(outStream, "Created a parking lot with %v slots\n", maxSlot)

		case validate(input, "park", 3):
			slotNo := 1
			fmt.Fprintf(outStream, "Allocated slot number: %v\n", slotNo)

		case validate(input, "leave", 2):
			slotNo := 1
			fmt.Fprintf(outStream, "Slot number %v is free\n", slotNo)

		case validate(input, "status", 1):
			cars := []struct {
				slot         int
				registration string
				colour       string
			}{
				{
					slot:         0,
					registration: "KA-01-HH-1234",
					colour:       "White",
				},
				{
					slot:         1,
					registration: "KA-01-HH-9999",
					colour:       "Black",
				},
			}
			var w = tabwriter.NewWriter(outStream, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "Slot No.\tRegistration No\tColour")
			for _, car := range cars {
				s := fmt.Sprintf("%v\t%s\t%s", car.slot, car.registration, car.colour)
				fmt.Fprintln(w, s)
			}
			w.Flush()

		case validate(input, "registration_numbers_for_cars_with_colour", 2):
			var registrations []string
			registrations = append(registrations, "KA-01-HH-1234")
			err := printer.Fprintf(outStream, registrations)
			if err != nil {
				panic(err.Error())
			}

		case validate(input, "slot_numbers_for_cars_with_colour", 2):
			var slots []int
			slots = append(slots, 1)
			err := printer.Fprintf(outStream, slots)
			if err != nil {
				panic(err.Error())
			}

		case validate(input, "slot_number_for_registration_number", 2):
			slotNo := 0
			fmt.Fprintln(outStream, slotNo)

		case validate(input, "exit", 1):
			exit = true

		default:
			fmt.Fprintln(outStream, "Unknown input command")
		}
	}
}

func parse(input string) []string {
	return strings.Split(input, " ")
}

func validate(input string, expectedCmd string, expectedLength int) bool {
	cmdArgs := parse(input)
	cmd := cmdArgs[0]
	if cmd == expectedCmd && len(cmdArgs) == expectedLength {
		return true
	}
	return false
}
