package cod

import (
	"github.com/adabei/goldenbot/rcon"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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
  events chan interface{}
}

func NewIntegrated(requests rcon.RCONQuery, ea events.Aggregator) *Integrated {
  i := new(Integrated)
  i.requests = requests
  i.events = ea.Subscribe(v)
  return i
}

func (i *Integrated) Setup() error {
  return nil
}

func (i *Integrated) Start() {
  //select{}
}

func Num(id string) int {
  ch := make(chan []byte)
  i.requests <- rcon.RCONQuery{Command: "status", Response: ch}
  res := <-ch
  if res != nil {
    status = string(res)
    players := strings.Split(status, "\n")
    for _, p := range players {
      sm := SubmatchMap(statusRegexp, statusRegexp.FindStringSubmatch(p))
      if sm != nil {
        if sm["guid"] == id {
          num, err := strconv.Atoi(sm["num"])
          if err != nil {
            log.Fatal(err)
          }
          return num
        }
      }
    }
  }

	return -1
}

// SubmatchMap returns a map of of named group matches.
// It returns nil if the regexp doesn't match.
// TODO Refractor
func SubmatchMap(re *regexp.Regexp, match []string) map[string]string {
	if match == nil {
		return nil
	}

	m := make(map[string]string)
	for i, v := range re.SubexpNames()[1:] {
		m[v] = match[i+1]
	}

	return m
}
