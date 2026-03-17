 #!/bin/bash

# Crear red
docker network create echo-test-net

# Levantar y correr el server, conectandolo a la red
docker build -t echo-server -f server/Dockerfile .
docker run -d --name echo-server --network echo-test-net echo-server -c "python main.py"

sleep 10

# Crear container que usa netcat (Cliente)
docker build -t echo-test -f echo-test/Dockerfile .
docker run --network echo-test-net echo-test -c "echo hola | nc echo-server 12345"

# Limpieza
docker rm -f echo-server
docker network rm echo-test-net
