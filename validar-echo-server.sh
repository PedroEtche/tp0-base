#!/bin/bash

MENSAJE="Docker testing..."

RESPUESTA=$(docker run --rm --network=tp0_testing_net alpine /bin/sh -c "echo '$MENSAJE' | nc server 12345")

# Validar
if [ "$RESPUESTA" = "$MENSAJE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
