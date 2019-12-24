package cmd

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func Test_cmd(t *testing.T) {
	// Save current state before rewriting state so that we can restore it later
	currArgs := os.Args
	currInputInteractive := inputInteractive
	currOutStream := outStream
	defer func() {
		// Restore state
		os.Args = currArgs
		inputInteractive = currInputInteractive
		outStream = currOutStream
	}()

	// Wire up interactive inputs redirection
	inpFile, err := os.Open("../test/input_interactive.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { inpFile.Close() }()
	inputInteractive = inpFile

	// Wire up outputs redirection
	var gotBuf bytes.Buffer
	outStream = &gotBuf

	// Expected CLI output
	out := `Created a parking lot with 6 slots
Allocated slot number: 1
Allocated slot number: 2
Allocated slot number: 3
Allocated slot number: 4
Allocated slot number: 5
Allocated slot number: 6
Slot number 4 is free
Slot No.    Registration No    Colour
1           KA-01-HH-1234      White
2           KA-01-HH-9999      White
3           KA-01-BB-0001      Black
5           KA-01-HH-2701      Blue
6           KA-01-HH-3141      Black
Allocated slot number: 4
Sorry, parking lot is full
KA-01-HH-1234, KA-01-HH-9999, KA-01-P-333
1, 2, 4
6
Not found
Not found
Not found
Unknown input command
`
	wantBuf := bytes.NewBufferString(out).Bytes()

	// Test cases
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "File input",
			args: []string{"cmd", "../test/input_file.txt"},
		},
		{
			name: "Interactive input",
			args: []string{"cmd"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			Run(os.Args)
			if !bytes.Equal(gotBuf.Bytes(), wantBuf) {
				t.Errorf("got = %v, want = %v", gotBuf.String(), string(wantBuf))
			}
		})
		gotBuf.Reset()
	}
}
