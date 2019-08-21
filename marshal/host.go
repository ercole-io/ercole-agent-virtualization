package marshal

import (
	"bufio"
	"encoding/json"
	"log"
	"strings"

	"github.com/ercole-io/ercole-agent-virtualization/model"
)

// Host returns a Host struct from the output of the host
// fetcher command. Host fields output is in key: value format separated by a newline
func Host(cmdOutput []byte) model.Host {

	lines := "{"
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	for scanner.Scan() {
		line := scanner.Text()
		splitted := strings.Split(line, ":")
		key := strings.Trim(splitted[0], " ")
		value := strings.Trim(splitted[1], " ")
		lines += marshalKey(key) + marshalValue(value) + ", "
	}

	lines += "}"
	lines = strings.Replace(lines, ", }", "}", -1)

	b := []byte(lines)
	var m model.Host
	err := json.Unmarshal(b, &m)

	if err != nil {
		log.Fatal(err)
	}

	return m
}
