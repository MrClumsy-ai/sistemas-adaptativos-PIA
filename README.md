# TODO

- [ ] Front-end: usar datos proporcionados por la API (ver documentacion en main.go)

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

4. Esperar a que en la terminal salga algo como:
```
WARNING:absl:Compiled the loaded model, but the compiled metrics have yet to be built. `model.compile_metrics` will be empty until you train or evaluate the model.
WARNING:absl:Compiled the loaded model, but the compiled metrics have yet to be built. `model.compile_metrics` will be empty until you train or evaluate the model.
WARNING:werkzeug: * Debugger is active!
INFO:werkzeug: * Debugger PIN: 256-911-240
```
5. Ir a la ubicacion desde un navegador web: http://localhost:8080/
6. Una vez que hicieron cambios, para verlos en accion, regresar a la terminal y presionar <kbd>Ctrl</kbd> + <kbd>C</kbd> y regresar al paso #3

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
git commit -m "mensaje corto describiendo tus cambios" -m "mensaje mas largo describiendo los cambios"
git pull --rebase
git push
```
15. Si es que sale algún error en
```git pull --rebase```
abortar el rebase usando
```git rebase --abort```
usar ahora un git pull regular y arreglar los conflictos de merge
```git pull```
una vez resueltos los conflictos, hacer otro commit con el mensaje de "merge"
```
git add .
git commit -m "merge"
git push
```
