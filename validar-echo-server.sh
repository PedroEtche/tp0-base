#!/bin/bash

MENSAJE="Docker testing..."

# Crear container que usa netcat (Cliente)
docker build -t echo-test -f echo-test/Dockerfile . 
# Correr cliente (netcat) y capturar respuesta 
RESPUESTA=$(docker run --rm --name echo-test --network tp0_testing_net echo-test \
    -c "echo $MENSAJE | nc server 12345")

# Validar
if [ "$RESPUESTA" = "$MENSAJE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
