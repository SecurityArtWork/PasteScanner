# PasteScanner
Pastescanner is a (Golang)tool to monitor websites content for relevant information

El programa no muestra nada por pantalla, pero va generado archivos en ./pastes con los pastes que no caducan, y en ./pastes/temp con los pastes con tiempo de vida.

Para pasarle los parametros editar paste.conf, tiene dos etiquetas bastente descriptivas, puedes ponerle todas las claves que quieras para que filtre con ellas y puedes quitar lineas de pastesitios para que no busque en ellos (añadir no, por que no he implementado más pero en el codigo verás que es copiar, pegar y editar dos tonterias en main() y find(), despues se puede añadir al .conf).
