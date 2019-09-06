package util

import (
	"fmt"
	"os/exec"
)

func OSExec(command string, args ...string) (output string, err error) {
	cmd := exec.Command(command, args...)
	co, e := cmd.CombinedOutput()
	//if err != nil {
	//	return "", err
	//}
	//if err = cmd.Start(); err != nil {
	//	return "", err
	//}
	////	out, _ := ioutil.ReadAll(co)
	////	e, _ := ioutil.ReadAll(se)
	//err = cmd.Wait()
	if len(co) > 0 {
		output = string(co)
	}

	if e != nil {
		if len(output) > 0 {
			err = fmt.Errorf(string(output))
		} else {
			err = e
		}
	}

	return
}

func bytesToStrings(out, e, co []uint8) (stdout, stderr, combined string) {
	stdout = ""
	stderr = ""
	if len(out) > 0 {
		stdout = string(out)
	}
	if len(e) > 0 {
		stderr = string(e)
	}
	if len(co) > 0 {
		combined = string(co)
	}
	return
}
