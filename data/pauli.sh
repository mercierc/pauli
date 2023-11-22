#!/bin/bash

# User defined functions.

function preinstall(){
    echo "preinstall not implemented."
}

function build(){
    echo "build not implemented."
}

function run(){
    echo "run not implemented"
}

function clean(){
    echo "clean not implemented."
}

function lint(){
    echo "lint not implemented."
}

function unittests(){
    echo "unittests not implemented."
}

function inttests(){
    echo "inttests not implemented."
}

# Common functions
function fatal(){
	echo -e "\e[1;4;97;48;5;9mFATAL: $1 \e[0m"
}
function highlight(){
	echo -e "\n\e[1;92m $1 \e[0m\n"
}

function info(){
	echo -e "\e[92m $1 \e[0m"
}

function warn(){
	echo -e "\e[38;5;208m $1 \e[0m"
}

preinstall

case $1 in
	build)
		build ;;
	run)
	        run ;;
	unittests)
		unittests "${@:2}" ;;
	lint)
		lint ;;
	inttests)
		inttests ;;
	staticanalysis)
		staticanalysis "${@:2}" ;;
	*)
		fatal "Unknown command"
		exit 1;;
esac
