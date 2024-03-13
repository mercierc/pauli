Pauli
=====

Description
-----------
The pauli project aims to assist and construct your application around a build image containing all the necessary tools and dependencies your project needs.
It simplify the use of a reproductible environment for developpement and integration by wrapping docker utilities that would be painfull to do with docker directly, such as: mounting volumes, select the right image, launch an interactice session...

Installation
-----------
To install the latest release version
```
# Install the latest release
VERSION=$(curl -I https://github.com/mercierc/pauli/releases/latest | grep -Eo 'v[0-9]+(\.[0-9]){2}')
wget https://github.com/mercierc/pauli/releases/download/$VERSION/pauli

sudo chown 700 pauli
sudo install pauli /usr/local/bin
```
To install a specific version set VERSION at vX.X.X, the desired version

Verify the installation with 
```
pauli --version
```

Quick start
---------
Complete the bash functions you need in the .pauli/pauli.sh file and call them inside the docker build container of your choice by invoking the pauli CLI. 

Let's build and run a simple go project. Create the following main.go file:
```
package main

import (
	"fmt"
	"os"
)

func main() {
        fmt.Printf("Hello, %s!", os.Getenv("WHO"))
}
```
In the same repository, initiate to the pauli project with the command:
`pauli init`
It will ask you for image name and tag and then create the _.pauli/_ folder. I choose an image of my choice, containing all the necessary to build and run my project. Here is the **.pauli/config.yaml** file created
```
builder:
  image: golang
  tag: alpine
  privileged: true
  volumes:
    - type: bind
      source: /var/run/docker.sock
      target: /var/run/docker.sock
name: go_example_with_pauli
```
Now, write the build and run functions according to your needs.
```
function build(){
		info build
		go build -o hello main.go
}

function run(){
    info run
    ./hello 
}
```
That's all! You can now build and run your executable in a build container undner a golang:alpine image:
`pauli build`

 and then run it:
`pauli run --env WHO=GO`

```
run 
Hello, GO!
```
Note that we can pass environement variables as we would do with docker.
