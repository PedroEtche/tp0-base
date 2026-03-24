### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.


Para resolver este ejercicio amplie el protocolo. Ahora el cliente debe mandar un avisando que va a mandar los batches. Una vez enviado esto puede hacer un mecanismo similar al del ejercicio anterior. Una vez terminado de mandar todo debe desconectarse. Ahora, para poder pedir por los ganadores, el cliente debe hacer un polling. Se conectara al servidor y mandara un mensaje pidiendo por los ganadores.

Estos dos nuevos mensajes (el de mandar batches y el de pedir ganadores) son de un solo byte y siempre que un cliente se conecta debe mandar este mensaje. El bit mas representativo indica el tipo de mensaje, y el resto de bit indica el ID del cliente. El ID es de mas utilidad cuando se piden los ganadores. Con esto nos evitamos hacer el broadcast