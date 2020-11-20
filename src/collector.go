/*  prometheus-postscreen-exporter
 *
 *  Copyright (C) 2020  Alexander Koch
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
  "time"
  "regexp"
  "strings"
  "log"
  "github.com/prometheus/client_golang/prometheus"
  "github.com/coreos/go-systemd/v22/sdjournal"
)

const (
  postfixUnit string = `postfix@-.service`
  pscrnPattern string = `^(PASS NEW|PASS OLD|PREGREET|DNSBL|HANGUP) `
)

var (
  counter = prometheus.NewCounterVec(
    prometheus.CounterOpts{
      Name: "postscreen_results",
      Help: "Number of postscreen results.",
    },
    []string{"result"},
  )
  pscrnRegex = regexp.MustCompile(pscrnPattern)
  journal *sdjournal.Journal
)

func readJournal() error {
  r := journal.Wait(time.Duration(5) * time.Second)
  if r < 0 {
    log.Fatal("Error waiting for journal")
  }

  for {
    c, err := journal.Next()
    if err != nil {
      return err
    }
    if c == 0 {
      break
    }

    e, err := journal.GetEntry()
    if err != nil {
      return err
    }

    matches := pscrnRegex.FindStringSubmatch(e.Fields["MESSAGE"])
    if matches != nil {
      result := strings.Replace(strings.ToLower(matches[1]), " ", "_", -1)
      counter.With(prometheus.Labels{"result":result}).Inc()
    }
  }

  return nil
}

func collect() {
  var err error

  journal, err = sdjournal.NewJournal()
  if err != nil {
    log.Fatal("Failed to open journal: " + err.Error())
  }
  err = journal.AddMatch("_SYSTEMD_UNIT=" + postfixUnit)
  if err != nil {
    log.Fatal("Failed to add unit filter: " + err.Error())
  }
  err = journal.SeekRealtimeUsec(uint64(time.Now().UnixNano() / 1000))
  if err != nil {
    log.Fatal("Failed to seek journal end: " + err.Error())
  }

  for {
    err = readJournal()
    if err != nil {
      log.Fatal("Failed to read journal: " + err.Error())
    }
  }
}

func init() {
  prometheus.MustRegister(counter)
  go collect()
}
