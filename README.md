### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.


#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).


### Respuesta

El protocolo que defini cumple los siguientes lineamientos:

El protocolo es binario y mixto. Tiene una parte que tamaño fija y otra variable. Los datos que se mandan son: 

- ID de la agencia: 1 byte -> unsigned int 8
- Largo del nombre: 1 byte -> unsigned int 8
- nombre: N bytes (acotado por el valor anterior entre 0 y 255): Cada byte representa un valor es ASCII
- Largo del apellido: 1 byte -> unsigned int 8
- apellido: N bytes (acotado por el valor anterior entre 0 y 255): Cada byte representa un valor es ASCII
- Documento (dni): 4 bytes -> unsigned int 32
- año de nacimiento: 2 bytes -> unsigned int 16
- mes de nacimiento: 1 byte -> unsigned int 8
- dia de nacimiento: 1 byte -> unsigned int 8
- Numero: 4 bytes -> unsigned int 32

Otros dos mensajes que existen son:

- ACK -> 1 byte
- NACK -> 1 byte

Cliente:

Como este primer ejercicio era mas de juguete que los siguientes 2, mantuve bastante el codigo provisto por la catedra.

Para comunicar una apuesta, el cliente lee de sus variables de entorno los datos. Crea el paquete que mencione anteriormente y se lo envia al servidor y espera una respuesta. Si el servidor valida la apuesta el cliente recibe un ACK, y si la apuesta es invalida se recibe un NACK

Servidor:

Aqui si cambie un poco mas el codigo.

En este caso, el servidor espera recibir el mensaje que comunica una apuesta apenas se conecta el cliente. Cuando la recibe la valida o invalida. Responde con ACK o NACK en cada caso

