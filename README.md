### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).


### Respuesta

Aqui utilice librerias que me permitan detectar `signals`. Opte por dejar que se termine la comunicacion si es que ya habia empezado. Tanto el cliente como el servidor tienen un flag que chequean para saber si llego o no un `SIGTERM`.

En el servidor, se configuró un timeout en el socket para aceptar nuevos clientes. Si se alcanza el timeout, el servidor verifica si recibió la señal `SIGTERM`. En caso de que la señal haya llegado, el servidor sale del loop de aceptación de clientes, cierra el socket y finaliza el programa. Si la señal `SIGTERM` se recibe mientras se está atendiendo a un cliente, el servidor finaliza la atención de ese cliente antes de verificar la señal. 

En el cliente la situacion es similar. Antes de conectarse al servivor vemos si llego la señal, si no lo hizo proseguimos. Si llego a la mitad de la conversacion, la misma prosigue hasta terminar, pero no se volvera a conectar si le faltaban mensajes por enviar.

Me parecio apropiado dejar que la conversacion termine dado que las mismas son muy cortas
