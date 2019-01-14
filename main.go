package main

import (
	"fmt"
	"time"
	"net/http"
	"encoding/json"
	"flag"
)
const (
	ENODE_FMT = "enode://%s@%s:%d"
	ADMIN_ADDPEER_FMT = `admin.addPeer("enode://%s@%s:%d")`
	NET_MAINNET		= 1
	NET_TESTNET		= 2
)

type NodeData struct {
	Id string		`json:"id"`
	Host string 	`json:"host"`
	Port int		`json:"port"`
	ClientId string `json:"clientId"`
	Client string 	`json:"client"`
	ClientVersion string 	`json:"clientVersion"`
	Os string 	`json:"os"`
	Country string `json:"country"`
}

type Result struct {
	Data []NodeData	`json:"data"`
}

func GetEthnodes( OutputFmt string, net, start, length int) ([]string, error) {

	res := make([]string,0)
	uri := fmt.Sprintf("https://www.ethernodes.org/network/%d/data?draw=1&columns%%5B0%%5D%%5Bdata%%5D=id&columns%%5B0%%5D%%5Bname%%5D=&columns%%5B0%%5D%%5Bsearchable%%5D=true&columns%%5B0%%5D%%5Borderable%%5D=true&columns%%5B0%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B0%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B1%%5D%%5Bdata%%5D=host&columns%%5B1%%5D%%5Bname%%5D=&columns%%5B1%%5D%%5Bsearchable%%5D=true&columns%%5B1%%5D%%5Borderable%%5D=true&columns%%5B1%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B1%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B2%%5D%%5Bdata%%5D=port&columns%%5B2%%5D%%5Bname%%5D=&columns%%5B2%%5D%%5Bsearchable%%5D=true&columns%%5B2%%5D%%5Borderable%%5D=true&columns%%5B2%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B2%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B3%%5D%%5Bdata%%5D=country&columns%%5B3%%5D%%5Bname%%5D=&columns%%5B3%%5D%%5Bsearchable%%5D=true&columns%%5B3%%5D%%5Borderable%%5D=true&columns%%5B3%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B3%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B4%%5D%%5Bdata%%5D=clientId&columns%%5B4%%5D%%5Bname%%5D=&columns%%5B4%%5D%%5Bsearchable%%5D=true&columns%%5B4%%5D%%5Borderable%%5D=true&columns%%5B4%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B4%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B5%%5D%%5Bdata%%5D=client&columns%%5B5%%5D%%5Bname%%5D=&columns%%5B5%%5D%%5Bsearchable%%5D=true&columns%%5B5%%5D%%5Borderable%%5D=true&columns%%5B5%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B5%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B6%%5D%%5Bdata%%5D=clientVersion&columns%%5B6%%5D%%5Bname%%5D=&columns%%5B6%%5D%%5Bsearchable%%5D=true&columns%%5B6%%5D%%5Borderable%%5D=true&columns%%5B6%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B6%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B7%%5D%%5Bdata%%5D=os&columns%%5B7%%5D%%5Bname%%5D=&columns%%5B7%%5D%%5Bsearchable%%5D=true&columns%%5B7%%5D%%5Borderable%%5D=true&columns%%5B7%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B7%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&columns%%5B8%%5D%%5Bdata%%5D=lastUpdate&columns%%5B8%%5D%%5Bname%%5D=&columns%%5B8%%5D%%5Bsearchable%%5D=true&columns%%5B8%%5D%%5Borderable%%5D=true&columns%%5B8%%5D%%5Bsearch%%5D%%5Bvalue%%5D=&columns%%5B8%%5D%%5Bsearch%%5D%%5Bregex%%5D=false&order%%5B0%%5D%%5Bcolumn%%5D=3&order%%5B0%%5D%%5Bdir%%5D=asc&start=%d&length=%d&search%%5Bvalue%%5D=30303&search%%5Bregex%%5D=false&_=%d", net, start, length, time.Now().UnixNano()/1e6)
	resp, err := http.Get(uri)
	if err != nil {
		return res, err
	}

	// Close the response once we return from the function.
	defer resp.Body.Close()

	// Check the status code for a 200 so we know we have received a
	// proper response.
	if resp.StatusCode != 200 {
		return res, fmt.Errorf("HTTP Response Error %d\n", resp.StatusCode)
	}

	// Decode the rss feed document into our struct type.
	// We don't need to check for errors, the caller can do this.
	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return res, err
	}
	if len(result.Data) == 0 {
		return res, fmt.Errorf("Error data format of result")
	}

	for _,x := range result.Data {
		res = append(res, fmt.Sprintf(OutputFmt, x.Id, x.Host, x.Port))
	}

	return res, nil
}


var flagStart int
var flagLength int
var flagFmt string
var flagNet string
func init() {
	flag.StringVar(&flagNet, "net", "mainnet", "mainnet or testnet")
	flag.IntVar(&flagStart, "start", 0, "start index")
	flag.IntVar(&flagLength, "length", 10, "length ")
	flag.StringVar(&flagFmt, "fmt", "addpeer", "enode or addpeer")
}

func main(){
	flag.Parse()

	fmtString := ADMIN_ADDPEER_FMT
	if flagFmt == "enode" {
		fmtString = ENODE_FMT
	}

	net := NET_MAINNET
	if flagNet == "testnet" {
		net = NET_TESTNET
	}

	l, err := GetEthnodes(fmtString, net, flagStart, flagLength)
	if err != nil {
		fmt.Printf("%v",err)
	}

	for _, e := range l {
		fmt.Println(e)
	}

}


