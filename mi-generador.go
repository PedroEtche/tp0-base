package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var docker_compose_server_config string = `name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net

`

var docker_compose_network_config string = `networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
`

func main() {
	outPutFile := os.Args[1]
	clientsLen := os.Args[2]

	clientsInt, err := strconv.Atoi(clientsLen)
	if err != nil {
		log.Fatalf("Invalid argument: %v", err)
	}

	f, err := os.Create(outPutFile)
	if err != nil {
		log.Fatalf("Error openning the file (maybe invalid argument was received): %v", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	writeConfig(w, docker_compose_server_config)

	writeClientsConfig(clientsInt, w)

	writeConfig(w, docker_compose_network_config)
}

func writeClientsConfig(clientsInt int, w *bufio.Writer) {
	for clientID := 1; clientID <= clientsInt; clientID++ {
		clientService := "client" + strconv.Itoa(clientID)
		cliID := strconv.Itoa(clientID)
		block := fmt.Sprintf(`  %s:
    container_name: %s
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=%s
      - CLI_NAME_FIRST=Santiago Lionel
      - CLI_NAME_LAST=Lorca
      - CLI_DOCUMENT=30904465
      - CLI_BIRTH_YEAR=1999
      - CLI_BIRTH_MONTH=03
      - CLI_BIRTH_DAY=17
      - CLI_NUMBER=7574
    volumes:
      - ./client/config.yaml:/config.yaml
    networks:
      - testing_net
    depends_on:
      - server

`, clientService, clientService, cliID)

		writeConfig(w, block)
	}
}

func writeConfig(w *bufio.Writer, s string) {
	defer w.Flush()
	// NOTE: Buffio parece evitar o al menos hacer muy poco probable que haya un shor write (al escribir en disco por lo menos).
	n, err := w.WriteString(s)
	if err != nil {
		log.Fatalf("Error appear when writting to file: %v", err)
	}
	fmt.Printf("wrote %d bytes\n", n)
}
