### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).


### Respuesta

Para resolver el ejercicio, realicé los siguientes cambios:

1. Modificación del docker-compose: Cambié la definición del archivo docker-compose para que tanto el servidor como los clientes utilicen `docker volumes`. Esto permite persistir los archivos de configuración fuera de las imágenes de Docker, asegurando que los cambios realizados en estos archivos sean efectivos sin necesidad de reconstruir las imágenes. Este cambio está reflejado en el script `mi-generador.go`.

2. Actualización de los `Dockerfile`:

    - Cliente: Comenté la línea que copiaba el archivo de configuración al interior de la imagen.
    - Servidor: Opté por seguir copiando todo el contenido del directorio dentro de la imagen, pero agregué un archivo `.dockerignore` para excluir el archivo de configuración durante el proceso de creación de la imagen.