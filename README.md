# Fully-Synched Cross-Platform To-Do List

A simple, quickly put-together example of a full-stack application using Golang, MongoDB and React Native. The mobile client uses React Native to deploy on both iOS and Android while the REST API is written in Golang and uses MongoDB for persistence. 


## Backend

The backend exposes a Golang REST API that is used when the client updates the list. The list is stored in a MongoDB database using basic CRUD commands from the Golang server. 

To run, first ensure that both golang 1.10+ and MongoDB are installed. Then open a terminal and start a MongoDB database by running `mongod` on MacOS or the equivalent on your own OS.Then change directory to the main folder of this repository and use `go run main.go` to start running the Golang server. 


## Frontend

The frontend uses React Native to enable both iOS and Android versions of the to-do list to be built. The client app fetchs and updates the list after each item addition, deletion or update.

To run, first make sure that react native and either the ios simulator or an android emulator are downloaded and setup. Then in the reactnative directory, run either `react-native run-android` or `react-native run-ios`.

