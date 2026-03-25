### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.


### Respuesta

Cambie el protocolo del ejercicio anterior:

Ahora hay un nuevo mensaje (mensaje `batch`):

- Largo del batch: 2 bytes -> unsigned int 16
- Payload: Mensaje del ejercicio anterior

Luego, mantuve los mensajes de ACK y NACK

Ya con este cambio el servidor puede entender correntamente los batches. Solo tiene que leer cuantas puestas hay y listo. En el caso de que el batcj este correcto se responde con un ACK y el cliente sigue enviando batches. Si el batch es invalido el servidor responde con NACK y el cliente descarta el batch y envia los siguientes.

Cuando el cliente manda todo cierra la conexion. El servidor toma este cierre como que todos los batches fueron mandados

En caso de que el batch que eliga la agencia para mandar sea muy grande, a la hora de armar el mensaje, se arma el paquete con un tamaño menor o igual a 8 kB. Si la agencia quiso armar un paquete que superaba este limite, debera guardarse las apuestas que no entraron en el batch y intentar mandarlas en el siguiente batch. Para hacer este manejo tenemos la funcion `CreateBatch(bets []Bet) ([]byte, []Bet)`. Ya en la misma firma de la funcion se puede ver que se devuelve un slice de `Bet`, esto es por si no entraron en el paquete y deberan ser mandadas en otro batch

```mermaid
sequenceDiagram
    participant Cliente
    participant Servidor

    Cliente->>Servidor: Conexión TCP abierta
    Cliente->>Servidor: Enviar Batch
    Servidor-->>Cliente: ACK
    Cliente->>Servidor: Enviar siguiente Batch
    Servidor-->>Cliente: ACK
    Cliente->>Servidor: Cerrar conexión TCP
    Cliente->>Cliente: Procesamiento completo