package impl

import (
	"fmt"
	// "github.com/alexbrainman/printer"
	"net"
	"os"
)

type LabelPrinter struct {
	LabelData string `json:"label_data"`
	Ip        string `json:"ip"`
	Port      string `json:"port"`
}

func printLabel(labelPrinter LabelPrinter) {
	// fmt.Printf("printLabel begin, label_data = %s, ip = %s, port = %s\n", labelPrinter.LabelData, labelPrinter.Ip, labelPrinter.Port)
	Info.Printf("printLabel begin, label_data = %s, ip = %s, port = %s", labelPrinter.LabelData, labelPrinter.Ip, labelPrinter.Port)
	// address -> "172.26.100.15:9100"
	if len(labelPrinter.Ip) == 0 {
		// TODO print by default local printer
		// printLabelByLocalPrinter(labelPrinter)
		return
	}
	address := labelPrinter.Ip + ":" + labelPrinter.Port
	// fmt.Printf("printLabel address = %s\n", address)
	Info.Printf("printLabel address = %s", address)
	conn, err := net.Dial("tcp", address)
	// checkError(err)

	if err != nil {
		// fmt.Printf("printLabel net.Dial happen error: %s", err.Error())
		Error.Printf("printLabel net.Dial happen error: %s", err.Error())
		return
	}

	// _, err = conn.Write([]byte("^XA^LH55,30^FO20,10^CFD,27,13^FDCompany Name^FS^FO20,60^AD^FDTESTDESC^FS^FO40,160^BY2,2.0^BCN,100,Y,N,N,N^FDTEST^FS^XZ"))
	_, err = conn.Write([]byte(labelPrinter.LabelData))
	// checkError(err)
	if err != nil {
		// fmt.Printf("printLabel conn.Write error: %s", err.Error())
		Error.Printf("printLabel conn.Write happen error: %s", err.Error())
		return
	}
	// fmt.Println("Print ZPL success!")
	Info.Printf("Print ZPL success!")

	// os.Exit(0)
	err = conn.Close()
	if err != nil {
		// fmt.Printf("printLabel conn.Close error: %s", err.Error())
		Error.Printf("printLabel conn.Close happen error: %s", err.Error())
		return
	}
}

func checkError(err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
