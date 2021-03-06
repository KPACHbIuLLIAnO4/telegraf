// +build darwin

package net

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/common"
)

func NetIOCounters(pernic bool) ([]NetIOCountersStat, error) {
	out, err := exec.Command("/usr/sbin/netstat", "-ibdn").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	ret := make([]NetIOCountersStat, 0, len(lines)-1)
	exists := make([]string, 0, len(ret))

	for _, line := range lines {
		values := strings.Fields(line)
		if len(values) < 1 || values[0] == "Name" {
			// skip first line
			continue
		}
		if common.StringsHas(exists, values[0]) {
			// skip if already get
			continue
		}
		exists = append(exists, values[0])

		base := 1
		// sometimes Address is ommitted
		if len(values) < 11 {
			base = 0
		}

		parsed := make([]uint64, 0, 6)
		vv := []string{
			values[base+3], // Ipkts == PacketsRecv
			values[base+4], // Ierrs == Errin
			values[base+5], // Ibytes == BytesRecv
			values[base+6], // Opkts == PacketsSent
			values[base+7], // Oerrs == Errout
			values[base+8], // Obytes == BytesSent
		}
		for _, target := range vv {
			if target == "-" {
				parsed = append(parsed, 0)
				continue
			}

			t, err := strconv.ParseUint(target, 10, 64)
			if err != nil {
				return nil, err
			}
			parsed = append(parsed, t)
		}

		n := NetIOCountersStat{
			Name:        values[0],
			PacketsRecv: parsed[0],
			Errin:       parsed[1],
			BytesRecv:   parsed[2],
			PacketsSent: parsed[3],
			Errout:      parsed[4],
			BytesSent:   parsed[5],
		}
		ret = append(ret, n)
	}

	if pernic == false {
		return getNetIOCountersAll(ret)
	}

	return ret, nil
}
