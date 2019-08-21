package marshal

import (
	"bufio"
	"encoding/json"
	"log"
	"strings"

	"github.com/ercole-io/ercole-agent-virtualization/model"
)

// Filesystems returns a list of Filesystem entries extracted
// from the filesystem fetcher command output.
// Filesystem output is a list of filesystem entries with positional attribute columns
// separated by one or more spaces
func Filesystems(cmdOutput []byte) []model.Filesystem {

	lines := "["
	scanner := bufio.NewScanner(strings.NewReader(string(cmdOutput)))
	for scanner.Scan() {
		lines += "{"
		line := scanner.Text()
		line = strings.Join(strings.Fields(line), " ")
		splitted := strings.Split(line, " ")
		lines += marshalKey("filesystem") + marshalString(strings.Trim(splitted[0], " ")) + ", "
		lines += marshalKey("fstype") + marshalString(strings.Trim(splitted[1], " ")) + ", "
		lines += marshalKey("size") + marshalString(strings.Trim(splitted[2], " ")) + ", "
		lines += marshalKey("used") + marshalString(strings.Trim(splitted[3], " ")) + ", "
		lines += marshalKey("available") + marshalString(strings.Trim(splitted[4], " ")) + ", "
		lines += marshalKey("usedperc") + marshalString(strings.Trim(splitted[5], " ")) + ", "
		lines += marshalKey("mountedon") + marshalString(strings.Trim(splitted[6], " ")) + ", "
		lines += "},"
	}

	lines += "]"
	lines = strings.Replace(lines, ", }", "}", -1)
	lines = strings.Replace(lines, "},]", "}]", -1)

	b := []byte(lines)
	var m []model.Filesystem
	err := json.Unmarshal(b, &m)

	if err != nil {
		log.Fatal(err)
	}

	return m
}
