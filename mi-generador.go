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
	n, err := w.WriteString(s)
	if err != nil {
		log.Fatalf("Error appear when writting to file: %v", err)
	}
	// TODO: Check for short write
	fmt.Printf("wrote %d bytes\n", n)
	w.Flush()
}
