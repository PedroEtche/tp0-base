### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.


Para este ejercicio modidfique el protocolo del ejercicio 5. Ahora el servidor solo espera el mensaje `batch`. Este mensaje esta compuesto del mismo protocolo del ejercicio anterior, pero se le suma un header extra. Los primero 2 bytes (uint16) del paquete indican la cantidad de apuestas (`bet`) que hay en el mensaje. Con este nuevo campo, sumado al protocolo del ejercicio 5, el servidor ya puede deserializar los batches correctamente. 

En el caso de que el batch este correcto el servidor responde con un `ack` al cliente. Cuando el cliente recibe el `ack` arma otro batch y lo manda. Este bucle se repite hasta que no haya mas apuestas. El cliente cierra la conexion cuando recibe el ultimo `ack` y no tiene mas apuestas. El servidor detecta el cierre y da por terminada la comunicacion.

En el caso de que haya algun problema en el batch, el servidor responde con un `nack` y descarta todo el batch. El cliente recibe el `nack` y lo imprime por pantalla. En este caso se descarta el batch y se continua con el siguiente.

En caso de que el batch que eliga la agencia para mandar sea muy grande, a la hora de armar el mensaje, se arma el paquete con un tamaño menor o igual a 8 kB. Si la agencia quiso armar un paquete que superaba este limite, debera guardarse las apuestas que no entraron en el batch y intentar mandarlas en el siguiente batch. Para hacer este manejo tenemos la funcion `CreateBatch(bets []Bet) ([]byte, []Bet)`. Ya en la misma firma de la funcion se puede ver que se devuelve un slice de `Bet`, esto es por si no entraron en el paquete y deberan ser mandadas en otro batch