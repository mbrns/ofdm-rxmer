package query

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gosnmp/gosnmp"
)

const (
	docsPnmBulkDestIpAddrTypeOID = "1.3.6.1.4.1.4491.2.1.27.1.1.1.1.0"
	docsPnmBulkDestIpAddrOID     = "1.3.6.1.4.1.4491.2.1.27.1.1.1.2.0"
	docsPnmBulkDestPathOID       = "1.3.6.1.4.1.4491.2.1.27.1.1.1.3.0"
	docsPnmBulkUploadControlOID  = "1.3.6.1.4.1.4491.2.1.27.1.1.1.4.0"
)

var (
	ofdmIfIndex                       = -1
	cmMacIfIndex                      = -1
	docsPnmCmDsOfdmRxMerFileEnableOID = "1.3.6.1.4.1.4491.2.1.27.1.2.5.1.1"
	docsPnmCmDsOfdmRxMerFileNameOID   = ".1.3.6.1.4.1.4491.2.1.27.1.2.5.1.8"
)

func getOfdmIndex(cmIp string, writeCommunity string) {
	gosnmp.Default.Target = cmIp
	gosnmp.Default.Community = writeCommunity
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v\n", err)
	}
	defer gosnmp.Default.Conn.Close()

	ifType := "1.3.6.1.2.1.2.2.1.3"

	err = gosnmp.Default.BulkWalk(ifType, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
}

func getCmMacAddress(cmIp string, writeCommunity string) string {
	gosnmp.Default.Target = cmIp
	gosnmp.Default.Community = writeCommunity
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v\n", err)
	}
	defer gosnmp.Default.Conn.Close()

	ifPhysAddressOID := fmt.Sprintf("%s.%s", "1.3.6.1.2.1.2.2.1.6", strconv.Itoa(cmMacIfIndex))

	res, err := gosnmp.Default.Get([]string{ifPhysAddressOID})
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	m := res.Variables[0].Value.([]uint8)

	return fmt.Sprintf("%x%x%x%x%x%x",
		m[0], m[1], m[2], m[3], m[4], m[5])

}

func printValue(pdu gosnmp.SnmpPDU) error {
	switch pdu.Type {
	case gosnmp.Integer:
		r := pdu.Value.(int)
		if r == 277 {
			oid := strings.Split(pdu.Name, ".")
			idx, _ := strconv.Atoi(oid[len(oid)-1])
			log.Printf("Found OFDM ifIndex at %d\n", idx)
			ofdmIfIndex = idx
		}
		if r == 127 {
			oid := strings.Split(pdu.Name, ".")
			idx, _ := strconv.Atoi(oid[len(oid)-1])
			log.Printf("Found CM MAC layer ifIndex at %d\n", idx)
			cmMacIfIndex = idx
		}
	default:
	}
	return nil
}

func SendSet(cmIp string, writeCommunity string, tftpSrv string, tftpPath string) string {
	gosnmp.Default.Target = cmIp
	gosnmp.Default.Community = writeCommunity
	err := gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v\n", err)
	}
	defer gosnmp.Default.Conn.Close()

	getOfdmIndex(cmIp, writeCommunity)
	macAddr := getCmMacAddress(cmIp, writeCommunity)

	pdus := make([]gosnmp.SnmpPDU, 0)

	docsPnmBulkDestIpAddrType := gosnmp.SnmpPDU{
		Name:  docsPnmBulkDestIpAddrTypeOID,
		Type:  gosnmp.Integer,
		Value: 1, // 1 = IPv4
	}
	pdus = append(pdus, docsPnmBulkDestIpAddrType)

	docsPnmBulkDestIpAddr := gosnmp.SnmpPDU{
		Name:  docsPnmBulkDestIpAddrOID,
		Type:  gosnmp.OctetString,
		Value: tftpServerIpAddressFix(tftpSrv),
	}
	pdus = append(pdus, docsPnmBulkDestIpAddr)

	docsPnmBulkDestPath := gosnmp.SnmpPDU{
		Name:  docsPnmBulkDestPathOID,
		Type:  gosnmp.OctetString,
		Value: tftpPath,
	}
	pdus = append(pdus, docsPnmBulkDestPath)

	docsPnmBulkUploadControl := gosnmp.SnmpPDU{
		Name:  docsPnmBulkUploadControlOID,
		Type:  gosnmp.Integer,
		Value: 3, // autoUpload(3)
	}
	pdus = append(pdus, docsPnmBulkUploadControl)

	if ofdmIfIndex == -1 {
		panic(errors.New("couldn't get CM OFDM ifIndex"))
	}

	docsPnmCmDsOfdmRxMerFileName := gosnmp.SnmpPDU{
		Name:  docsPnmCmDsOfdmRxMerFileNameOID + "." + strconv.Itoa(ofdmIfIndex),
		Type:  gosnmp.OctetString,
		Value: macAddr + ".PNMDsMer",
	}
	pdus = append(pdus, docsPnmCmDsOfdmRxMerFileName)

	err = gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v\n", err)
	}
	defer gosnmp.Default.Conn.Close()

	_, err = gosnmp.Default.Set(pdus)

	if err != nil {
		log.Panicln(err)
	}

	err = gosnmp.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v\n", err)
	}
	defer gosnmp.Default.Conn.Close()

	docsPnmCmDsOfdmRxMerFileEnable := gosnmp.SnmpPDU{
		Name:  docsPnmCmDsOfdmRxMerFileEnableOID + "." + strconv.Itoa(ofdmIfIndex),
		Type:  gosnmp.Integer,
		Value: 1,
	}

	_, err = gosnmp.Default.Set([]gosnmp.SnmpPDU{docsPnmCmDsOfdmRxMerFileEnable})

	if err != nil {
		log.Panicln(err)
	}

	return tftpPath + macAddr + ".PNMDsMer"
}

func tftpServerIpAddressFix(ip string) []byte {
	i := strings.Split(ip, ".")
	oct1, _ := strconv.Atoi(i[0])
	oct2, _ := strconv.Atoi(i[1])
	oct3, _ := strconv.Atoi(i[2])
	oct4, _ := strconv.Atoi(i[3])
	return []byte{
		byte(oct1),
		byte(oct2),
		byte(oct3),
		byte(oct4),
	}
}
