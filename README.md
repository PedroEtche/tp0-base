### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.


### Respuesta

Cambie el protocolo del ejercicio anterior:

Ahora hay un nuevo mensaje que se usa para indicar que quiere hacer el cliente

- Tipo de request: 1 byte -> El bit mas significativo indica si el request va a ser en enviar batches o conultar por ganadores. El resto de los bits se usa para indicar el numero de la agencia
- Payload (en caso de queres cargar batches): Mensaje del ejercicio anterior

Mensaje para comunicar los ganadores:

- Ganadores: 1 byte -> Si el byte es 0xFF se toma como NACK y por lo tanto el cliente debe esperar. En caso contrario indica la cantidad de documentos de ganadores que hay en el mensaje
- Documentos: 4 bytes (unsigned int 32) -> Lo mismo que el ej5 


Cliente:

1- Envio de batches:
El cliente primero envía un mensaje para especificar que enviará un batch de apuestas.Luego, continúa la comunicación como en el ejercicio anterior (cerrando la conexión al finalizar el envío de apuestas).

2- Consulta de ganadores:
El cliente realiza polling para verificar si el sorteo ya se realizó. Al conectarse al servidor, envía el nuevo mensaje indicando que desea consultar por los ganadores.
Si recibe un NACK, cierra la conexión y espera 5 segundos antes de intentar nuevamente.
Si recibe la lista de ganadores, los procesa y los anuncia.

Servidor:
El mayor cambio es que se guarda la cantidad de agencias que hay. De esta forma sabe si puede iniciar o no el sorteo. Modifique el docker compose para que el `SERVER_LISTEN_BACKLOG` del server sea igual a la cantidad de clientes que existiran.

Detalle de implementacion: Para diferenciar la respuesta de ganadadores a la de NACK. El mensaje de ganadores tiene como primer byte la cantidad de ganadores. Por lo tanto, el cliente siempre lee un solo byte primero. Un NACK se representa con un byte igual a `0xFF`. Esto limita la cantidad máxima de ganadores a 254, lo cual asumi suficiente para este caso.


```mermaid
sequenceDiagram
    participant Cliente1
    participant Servidor
    participant Cliente2

    Cliente1->>Servidor: Conexión TCP abierta
    Cliente1->>Servidor: Voy a enviar Batch
    Cliente1->>Servidor: Enviar Batch
    Servidor-->>Cliente1: ACK
    Cliente1->>Servidor: Cerrar conexión TCP

    Cliente1->>Servidor: Conexión TCP abierta
    Cliente1->>Servidor: Ganadores(ID=1)
    Servidor-->>Cliente1: NACK
    Cliente1->>Servidor: Cerrar conexión TCP

    Cliente2->>Servidor: Conexión TCP abierta
    Cliente2->>Servidor: Voy a enviar Batch
    Cliente2->>Servidor: Enviar Batch
    Servidor-->>Cliente2: ACK
    Cliente2->>Servidor: Cerrar conexión TCP

    Cliente2->>Servidor: Conexión TCP abierta
    Cliente2->>Servidor: Ganadores(ID=2)
    Servidor-->>Cliente2: Lista de ganadores
    Cliente2->>Servidor: Cerrar conexión TCP

    Cliente1->>Servidor: Conexión TCP abierta
    Cliente1->>Servidor: Ganadores(ID=1)
    Servidor-->>Cliente1: Lista de ganadores
    Cliente1->>Servidor: Cerrar conexión TCP
