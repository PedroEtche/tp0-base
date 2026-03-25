## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).


### Respuesta


Para este ejercicio opte por no cambiar el protocolo y hacer solo cambios en el server. Dicho esto, encontre dos secciones criticas

1.  La primera se encuentra con las funciones provistas por la catedra

2.  La segunda con el contador que use en el ejercicio 7 para saber cuando ya se habian conectado todos los clientes

Ambos problemas los solucione con un `Lock`. Deje bien explicito en el codigo cuando empieza y termina la seccion critica con un comentario.
