### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.


### Respuesta

Se puede ejecutar el script de la misma forma que propone el ejercicio. La estructura del script es muy similar a la que aparece en el enunciado, con la diferencia de que utilicé Go en lugar de Python. Elegí Go porque es un lenguaje que me interesa aprender.

Para generar el archivo de Docker Compose, no utilicé ningún paquete especial para el manejo de archivos YAML. En su lugar, escribí directamente la definición en texto plano. Esto incluye las configuraciones del servidor, los clientes y la red.

La definición de los clientes está implementada dentro de un bucle, que se ejecuta tantas veces como se indique en el parámetro del script. Cada cliente tiene un número único que se incrementa automáticamente en cada iteración. Si el parámetro recibido es 0, entonces no se generan clientes.
