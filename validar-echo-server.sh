 #!/bin/bash

MENSAJE="Docker testing..."

# Levantar servicios
docker compose -f docker-compose-dev.yaml up -d --build

# Dar tiempo a que el servidor se levante
sleep 5

# Crear container que usa netcat (Cliente)
docker build -t echo-test -f echo-test/Dockerfile . 
# Correr cliente (netcat) y capturar respuesta 
RESPUESTA=$(docker run --name echo-test --network tp0_testing_net echo-test \
    -c "echo $MENSAJE | nc server 12345")

# Validar
if [ "$RESPUESTA" = "$MENSAJE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi

# Limpieza
docker compose -f docker-compose-dev.yaml stop -t 1
docker compose -f docker-compose-dev.yaml down
docker rm -f echo-test
