package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/mbrns/ofdm-rxmer/query"
	"github.com/pin/tftp"
)

func main() {
	// CM IP Address
	cmIpAddr := flag.String("cm", "", "Cable Modem IP address to poll for RxMER data.")

	// TFTP Server
	tftpIpAddr := flag.String("tftp", "", "TFTP server IP address to send rxMER data too.")

	// Write String
	snmpWriteString := flag.String("comm", "", "SNMP write string for cable modem.")

	// Select output
	outputType := flag.String("out", "json",
		"Specify output type for decoded DS PNM RxMER. Valid options are json, json-pretty, http-table, or http-chart")

	flag.Parse()

	if len(os.Args) < 5 {
		flag.Usage()
		os.Exit(1)
	}

	switch *outputType {
	case "json", "json-pretty", "http-chart", "http-table":
		break
	default:
		log.Println("ERROR INVALID OUTPUT TYPE")
		flag.Usage()
		os.Exit(1)
	}

	fname := query.SendSet(*cmIpAddr, *snmpWriteString, *tftpIpAddr)

	log.Println("Sleeping for 3 seconds to allow the TFTP transfer to complete")
	time.Sleep(3 * time.Second)

	c, err := tftp.NewClient(*tftpIpAddr + ":69")
	if err != nil {
		log.Fatalln(err)
	}
	wt, err := c.Receive(fname, "octet")
	if err != nil {
		log.Fatalln(err)
	}
	b := bytes.NewBuffer(nil)
	_, err = wt.WriteTo(b)
	if err != nil {
		log.Fatalln(err)
	}

	pnm := parseRxmerFile(b.Bytes())
	pnm.generateOutput(*outputType)
}

func (f PnmDsMerFile) generateOutput(out string) {
	tplData := make(map[Frequency]float64)
	for _, v := range f.RxmerData {
		tplData[v.Freq] = v.Mer
	}

	switch out {
	case "json":
		jsondata, err := json.Marshal(f)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(string(jsondata))
	case "json-pretty":
		jsondata, err := json.MarshalIndent(f, "", "    ")
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(string(jsondata))
	case "http-chart":
		t, err := template.New("tpl").Parse(chart)
		if err != nil {
			panic(err)
		}
		err = t.Execute(os.Stdout, tplData)
		if err != nil {
			panic(err)
		}
	case "http-table":
		t, err := template.New("tpl").Parse(table)
		if err != nil {
			panic(err)
		}
		err = t.Execute(os.Stdout, tplData)
		if err != nil {
			panic(err)
		}
	}

}
