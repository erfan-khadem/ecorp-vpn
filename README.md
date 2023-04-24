# README

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.
To build for a specific platform/arch, use `wails build -platform windows/amd64`


## Running
After building the project, download sing-box and put it next to the executable (and change the name to be exactly `sing-box` or `sing-box.exe`). Then in linux run the `build/bin/ecorp-linux-first-time.sh` script. After that you may start the app in Administrator mode in windows and normal user (without sudo) in linux.

## Screenshots
![Main Menu](https://github.com/er888kh/ecorp-vpn/blob/main/screenshots/main-menu.png?raw=true)
