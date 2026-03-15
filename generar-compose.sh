#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
go run mi-generador.go $1 $2
echo "Archivo $1 creado"
