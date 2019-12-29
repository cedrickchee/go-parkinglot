package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/cedrickchee/go-parkinglot/internal/printer"
)

type RunOptions struct {
	Stdin  io.Reader
	Stdout io.Writer
}

func Run(args []string) {
	RunCustom(args, nil)
}

func RunCustom(args []string, runOpts *RunOptions) {
	if runOpts == nil {
		runOpts = &RunOptions{}
	}

	if runOpts.Stdin == nil {
		runOpts.Stdin = os.Stdin
	}
	if runOpts.Stdout == nil {
		runOpts.Stdout = os.Stdout
	}

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
		scanner = bufio.NewScanner(runOpts.Stdin)
	}

	// Create a parking lot
	var parkinglot = &ParkingLot{}

	exit := false
	for !exit && scanner.Scan() {
		input := scanner.Text()

		cmdArgs := parse(input)

		switch {
		case validate(cmdArgs, "create_parking_lot", 2):
			capacity, err := strconv.Atoi(cmdArgs[1])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
				break
			}
			if err := parkinglot.createParkingLot("Marina Bay Sands", capacity); err == nil {
				fmt.Fprintf(runOpts.Stdout, "Created a parking lot with %v slots\n", capacity)
			} else {
				fmt.Fprintln(runOpts.Stdout, err.Error())
			}

		case validate(cmdArgs, "park", 3):
			slot, err := parkinglot.park(cmdArgs[1], cmdArgs[2])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
			} else {
				fmt.Fprintf(runOpts.Stdout, "Allocated slot number: %v\n", slot.getParkingSlotNumber())
			}

		case validate(cmdArgs, "leave", 2):
			slotNumber, err := strconv.Atoi(cmdArgs[1])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
				break
			}
			if err := parkinglot.leave(slotNumber); err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
			} else {
				fmt.Fprintf(runOpts.Stdout, "Slot number %v is free\n", slotNumber)
			}

		case validate(cmdArgs, "status", 1):
			slots := parkinglot.getStatus()
			var w = tabwriter.NewWriter(runOpts.Stdout, 0, 0, 4, ' ', 0)
			fmt.Fprintln(w, "Slot No.\tRegistration No\tColour")
			for _, slot := range slots {
				vehicle := slot.getVehicle()
				s := fmt.Sprintf("%v\t%s\t%s", slot.getParkingSlotNumber(), vehicle.getNumber(), vehicle.getColor())
				fmt.Fprintln(w, s)
			}
			w.Flush()

		case validate(cmdArgs, "registration_numbers_for_cars_with_colour", 2):
			_, regisNumbers, err := parkinglot.getVehiclesByColor(cmdArgs[1])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
				break
			}
			err = printer.Fprintf(runOpts.Stdout, regisNumbers)
			if err != nil {
				panic(err.Error())
			}

		case validate(cmdArgs, "slot_numbers_for_cars_with_colour", 2):
			slotNumbers, _, err := parkinglot.getVehiclesByColor(cmdArgs[1])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
				break
			}
			err = printer.Fprintf(runOpts.Stdout, slotNumbers)
			if err != nil {
				panic(err.Error())
			}

		case validate(cmdArgs, "slot_number_for_registration_number", 2):
			slotNumber, err := parkinglot.getVehicleByRegistrationNumber(cmdArgs[1])
			if err != nil {
				fmt.Fprintln(runOpts.Stdout, err.Error())
				break
			}
			fmt.Fprintln(runOpts.Stdout, slotNumber)

		case validate(cmdArgs, "exit", 1):
			exit = true

		default:
			fmt.Fprintln(runOpts.Stdout, "Unknown input command")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func parse(input string) []string {
	cutset := "\n"
	if runtime.GOOS == "windows" {
		cutset = "\r\n"
	}

	trimmed := strings.TrimRight(input, cutset)
	return strings.Split(trimmed, " ")
}

func validate(cmdArgs []string, expectedCmd string, expectedLength int) bool {
	cmd := cmdArgs[0]
	if cmd == expectedCmd && len(cmdArgs) == expectedLength {
		return true
	}
	return false
}
