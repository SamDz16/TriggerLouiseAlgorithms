# Trigger Louise Algorithms - Chrome/Firefox Extension
This is a chrome/firefow extensions which gets triggered whenever user heads to the `https://dbpedia.org/sparql` enpoint. It adds to the Document Object Model (DOM) of the actual page a checkbox that user can click and enable in order to change the default form submission. An express node js server is made in order to server the `fetch.wasm` file generated by compiling the `main.go` file inside the `wasm`folder.
isFailing is an algorithm which returns:
- 1: if there is a failing : no results were being fetched,
- 0: if there is a success (no failing) : at leeast one result has been fetched

The user has possibility to display or not (hide) the results that were fetched in the case where isFailing is false.

## Installation
In order to enable correctely this extension into Firefox, you have to have some prerequisites:
* Have Chrome/Firefox installed in your machine (quite obvious xD)
* Load your extension into your browser
* Go and install the GO programming language at https://go.dev (I am using the 1.16.12 version)
* Make sure to have node and npm installed
* Make sure to add the `C:\Program Files\Go\bin` to your PATH environment variable

## Step 1: Compile your GO file to WebAssembly
The bellowed command is executed in the root directory in order to compile the `RelaxAPI.go` to get the wasm equivalent `fetch.wasm` in the `./server/public` folder.

```shell
GOOS=js GOARCH=wasm go build -o ./server/public/fetch.wasm ./RelaxAPI.go
```

## Step 2: Install npm packages
In order to launch your node express server to serve your ressources, you have to have node modules installed inside your folder, for that, you have to execute this command:
```shell
npm i nodemon -g && cd server && npm install && npm start
```
Note: i've made an npm script called *start, it helps to quickly launch the server using nodemon

## Step 3: Open the endpoint
Open your Chrome/Firefox browser and head to this endpoint `https://dbpedia.org/sparql`. Enjoy!
