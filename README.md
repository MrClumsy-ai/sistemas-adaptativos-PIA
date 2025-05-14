## Para iniciar el servidor:

1. Instalar golang: https://go.dev/doc/install
2. Ir al directorio en una terminal e instalar todos los paquetes del proyecto

```
cd ubicacion/del/proyecto
go mod tidy
```

3. Correr el servidor con go run

```
go run main.go
```

4. Ir a la ubicacion desde un navegador web: http://localhost:8080/

# TODO

## Para trabajar

1. descargar git: https://git-scm.com/downloads
2. abrir la terminal y escribir:

```
git config --global user.name "UsuarioDeGit"
git config --global user.email usuariodegit@gmail.com
git config --global init.default branch main
```

3. ir a github > perfil > ajustes > developer settings > personal access tokens > tokens (classic)
4. generar un token nuevo (clasico)
5. agregar en la nota el nombre del dispositivo en el que vas a trabajar
6. ajustar la expiracion del token
7. seleccionar todas las opciones
8. generar el token y guardarlo en un lugar seguro (notepad o algun lugar que no se va a olvidar)
9. ir al repo para trabajar y copiar la url
10. ir al folder en donde se va a clonar el repo con

```
cd tu/ubicacion/
```

11. ir a la terminal y escribir:

```
git clone https://ghp_TuToken@github.com/Usuario/Repo
```

12. abrir con algun editor de texto (vs code, neovim, etc) con

```
code . # para vscode
nvim . # para neovim
```

13. hacer tus cambios al repo
14. para subir tus cambios, hacer:

```
git add .
git commit -m "mensaje corto describiendo tus cambios"
git push --rebase
```
