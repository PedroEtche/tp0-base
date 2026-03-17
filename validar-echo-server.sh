 #!/bin/bash

MENSAJE="Docker testing..."

# Crear red
docker network create echo-test-net

# Levantar y correr el server, conectandolo a la red
docker build -t echo-server -f server/Dockerfile .
docker run -d --name echo-server --network echo-test-net echo-server -c "python main.py"

# Dar tiempo a que el servidor se levante
sleep 5

# Crear container que usa netcat (Cliente)
docker build -t echo-test -f echo-test/Dockerfile .
# Correr cliente (netcat) y capturar respuesta 
RESPUESTA=$(docker run --name echo-test --network echo-test-net echo-test \
    -c "echo $MENSAJE | nc echo-server 12345")

# Validar
if [ "$RESPUESTA" = "$MENSAJE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi

# Limpieza
docker rm -f echo-server > /dev/null
docker rm -f echo-test > /dev/null
docker network rm echo-test-net > /dev/null
