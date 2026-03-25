### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `


### Respuesta

Para ejecutar el script y que no haya conflictos es preferible crear un compose sin cliente primero. Para eso se debe ejecutar `./generar-compose.sh docker-compose-dev.yaml 0`. Luego, se puede ejecutar es script `./validar-echo-server.sh`.

El script esta escrito totalmente en bash y es bastante sencillo. Se define un mensaje para hacer la prueba. Luego, se ejecuta lel siguiente comando: `docker run --rm --network=tp0_testing_net alpine /bin/sh -c "echo '$MENSAJE' | nc server 12345"`

Este comando corre un contenedor con una imagen pequeña que tiene la funcionalidad de `netcat`. El contenedor puede encontrar al `server` gracias a la resolucion de DNS que proporciona Docker. Es posible esto ya que con el flag `--network=tp0_testing_net` hacemos que los dos contenedores esten en la misma red.

Cuando termina de correr este container se captura el resultado y se verifica si cumple con la funcionalidad del echo server (responde lo mismo que le mandamos)
