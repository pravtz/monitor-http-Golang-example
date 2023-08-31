package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Server struct {
	ServerName  string
	ServerUrl   string
	Elapsed     float64
	StatusCode  int
	FailureDate string
}

func createServerList(serverList *os.File) []Server {
	csvReader := csv.NewReader(serverList)
	data, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	var servers []Server
	for i, line := range data {
		if i > 0 {
			server := Server{
				ServerName: line[0],
				ServerUrl:  line[1],
			}
			servers = append(servers, server)
		}

	}
	return servers
}

func checkServer(serversList []Server) []Server {
	var downServer []Server

	for _, server := range serversList {
		now := time.Now()
		get, err := http.Get(server.ServerUrl)
		if err != nil {
			fmt.Printf("server %s is down [%s]\n", server.FailureDate, err.Error())
			server.StatusCode = 0
			server.FailureDate = now.Format("02/01/2006 15:04:05")
			fmt.Println("An error occurred while executing get(url)")
			downServer = append(downServer, server)
			continue
		}
		server.StatusCode = get.StatusCode
		if server.StatusCode != 200 {
			server.FailureDate = now.Format("02/01/2006 15:04:05")
			downServer = append(downServer, server)
		}
		server.Elapsed = time.Since(now).Seconds()

		fmt.Printf("status: [%d] and slapsed is: [%f] secunds, [%s]: [%s]\n", server.StatusCode, server.Elapsed, server.ServerName, server.ServerUrl)

	}
	return downServer
}
func openFiles(serverListFile string, downTimeFile string) (*os.File, *os.File) {
	serverList, err := os.OpenFile(serverListFile, os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	downtimeList, err := os.OpenFile(downTimeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return serverList, downtimeList
}
func generateDowntime(downtimeList *os.File, downServer []Server) {
	csvWriter := csv.NewWriter(downtimeList)
	for _, servidor := range downServer {
		line := []string{servidor.ServerName, servidor.ServerUrl, servidor.FailureDate, fmt.Sprintf("%f", servidor.Elapsed), fmt.Sprintf("%d", servidor.StatusCode)}
		csvWriter.Write(line)
	}
	csvWriter.Flush()
}

func main() {
	serverList, downtimeList := openFiles(os.Args[1], os.Args[2])
	defer serverList.Close()
	defer downtimeList.Close()
	servers := createServerList(serverList)
	downServer := checkServer(servers)

	generateDowntime(downtimeList, downServer)

}
