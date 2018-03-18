package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type TCPFlowSample struct {
	ID          string  `json:"id"`
	Start       float64 `json:"start"`
	End         float64 `json:"stop"`
	Bytes       int     `json:"bytes"`
	Bps         float64 `json:"bps"`
	Retransmits int     `json:"retransmits"`
	SndCwnd     int     `json:"snd-cwnd"`
	RttMs       float64 `json:"rtt-ms"`
	RttVar      int     `json:"rtt-var"`
	Pmtu        int     `json:"pmtu"`
}

type UDPFlowSample struct {
	ID          string  `json:"id"`
	Start       float64 `json:"start"`
	End         float64 `json:"end"`
	Bytes       int     `json:"bytes"`
	Bps         float64 `json:"bps"`
	JitterMs    float64 `json:"jitter-ms"`
	LostPackets int     `json:"lost-packets"`
	LostPercent float64 `json:"lost-percent"`
	Packets     int     `json:"packets"`
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("missing (json) file")
	}

	raw, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	var tcpFlowStats TCPFlowStats
	var udpFlowStats UDPFlowStats

	err = json.Unmarshal([]byte(raw), &tcpFlowStats)
	if err == nil && tcpFlowStats.Start.TestStart.Protocol == "TCP" {
		crunchTCP(tcpFlowStats)
	} else if err = json.Unmarshal([]byte(raw), &udpFlowStats); err == nil {
		crunchUDP(udpFlowStats)
	} else {
		log.Fatal(err)
	}
}

func crunchTCP(tcpFlowStats TCPFlowStats) error {
	var tcpFlowSamples []TCPFlowSample
	var tcpFlowSample TCPFlowSample

	start := tcpFlowStats.Start.Timestamp.Timesecs

	for _, interval := range tcpFlowStats.Intervals {
		for _, stream := range interval.Streams {
			flowID := fmt.Sprintf("%s-%d-%d",
				tcpFlowStats.Title, start, stream.Socket)

			tcpFlowSample = TCPFlowSample{
				ID:          flowID,
				Start:       float64(start) + stream.Start,
				End:         float64(start) + stream.End,
				Bytes:       stream.Bytes,
				Bps:         stream.BitsPerSecond,
				Retransmits: stream.Retransmits,
				SndCwnd:     stream.SndCwnd,
				RttMs:       float64(stream.Rtt) / 1000,
				RttVar:      stream.Rttvar,
				Pmtu:        stream.Pmtu,
			}
		}

		tcpFlowSamples = append(tcpFlowSamples, tcpFlowSample)
	}

	out, err := json.Marshal(tcpFlowSamples)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}

func crunchUDP(udpFlowStats UDPFlowStats) error {
	var udpFlowSamples []UDPFlowSample
	var udpFlowSample UDPFlowSample

	start := udpFlowStats.ServerOutputJSON.Start.Timestamp.Timesecs

	for _, interval := range udpFlowStats.ServerOutputJSON.Intervals {
		for _, stream := range interval.Streams {
			flowID := fmt.Sprintf("%s-%d-%d",
				udpFlowStats.Title, start, stream.Socket)

			udpFlowSample = UDPFlowSample{
				ID:          flowID,
				Start:       float64(start) + stream.Start,
				End:         float64(start) + stream.End,
				Bytes:       stream.Bytes,
				Bps:         stream.BitsPerSecond,
				JitterMs:    stream.JitterMs,
				LostPackets: stream.LostPackets,
				LostPercent: stream.LostPercent,
				Packets:     stream.Packets,
			}
		}

		udpFlowSamples = append(udpFlowSamples, udpFlowSample)
	}

	out, err := json.Marshal(udpFlowSamples)
	if err != nil {
		return err
	}

	fmt.Println(string(out))

	return nil
}
