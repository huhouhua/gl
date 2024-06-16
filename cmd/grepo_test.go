// Copyright 2024 The huhouhua Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"
	"io"
	"os"
)

// cmdTestCase describes a test case that works with releases.
type cmdTestCase struct {
	name      string
	cmd       string
	wantError bool
}

func tInfo(msg interface{}) {
	fmt.Println("--- INFO:", msg)
}
func tOut(msg interface{}) {
	fmt.Println("--- OUTPUT:", msg)
}

// A helper to ignore os.Exit(1) errors when running a cobra Command
func executeCommand(cmd string) (stdout string, err error) {
	args, err := shellwords.Parse(cmd)
	if err != nil {
		return "", err
	}
	return executeCommandOfArgs(args...)
}

func executeCommandOfArgs(args ...string) (stdout string, err error) {
	buf := new(bytes.Buffer)
	root, err := NewRootCmd(buf)
	if err != nil {
		return "", err
	}
	for i, arg := range args {
		tInfo(fmt.Sprintf("(%d) %s", i, arg))
	}
	// programs can exit with error here..
	stdout, _, err = executeCommandC(root, args...)
	tInfo("The command successfully returned values for assertion.")
	return stdout, err

}

// A helper to ignore os.Exit(1) errors when running a cobra Command
func executeCommandC(root *cobra.Command, args ...string) (stdout string, output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	// see https://stackoverflow.com/questions/10473800/in-go-how-do-i-capture-stdout-of-a-function-into-a-string
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	_, err = root.ExecuteC()

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			panic(err)
		}
		outC <- buf.String()
	}()

	// back to normal state
	if err := w.Close(); err != nil {
		tInfo(err)
	}
	os.Stdout = old // restoring the real stdout
	stdout = <-outC

	return stdout, buf.String(), err
}
