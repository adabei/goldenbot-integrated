package cod

import (
	"github.com/adabei/goldenbot/rcon"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var status string
var statusRegexp *regexp.Regexp = regexp.MustCompile("(?P<num>[0-9]+)" +
	"\\s+(?P<score>-?[0-9]+)" +
	"\\s+(?P<ping>[0-9]+)" +
	"\\s+(?P<guid>[0-9a-f]+)" +
	"\\s+(?P<name>.*)" +
	"\\s+(?P<lastmsg>[0-9]+)" +
	"\\s+(?P<address>(?:[0-9]{1,3}\\.){3}[0-9{1,3}:[0-9]{1,5})" +
	"\\s*(?P<qport>[0-9]{1,5})" +
	"\\s+(?P<rate>[0-9]+)")

type Integrated struct {
	requests chan rcon.RCONQuery
}

func (i *Integrated) Start() {
	for {
		ch := make(chan []byte)
		i.requests <- rcon.RCONQuery{Command: "status", Response: ch}
		res := <-ch
		status = string(res)
		time.Sleep(10000 * time.Millisecond)
	}
}

func Num(id string) int {
	players := strings.Split(status, "\n")
	for _, p := range players[2:] {
		sm := SubmatchMap(statusRegexp, statusRegexp.FindStringSubmatch(p))
		if sm["guid"] == id {
			num, err := strconv.Atoi(sm["num"])
			if err != nil {
				log.Fatal(err)
			}
			return num
		}
	}

	return -1
}

func SubmatchMap(re *regexp.Regexp, match []string) map[string]string {
	m := make(map[string]string)
	for i, v := range re.SubexpNames()[1:] {
		m[v] = match[i+1]
	}

	return m
}
