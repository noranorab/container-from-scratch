package main

import (
	"time"
)

type container struct {
	id          string    //unique identifier for the container
	name        string    //a name for the container
	image       string    //reference to the image used for the container
	status      string    //"created", "running", "exited"
	command     string    //command to run within the container
	createdTime time.Time //Timestamp when the container was created
}

type ContainerConfig struct {
	image       string            //name of the container image
	volumes     map[string]string //volume mounts
	ports       map[int]int       //port mapping
	Environment map[string]string //environment variables
}

type Container interface {
	//Start the container
	start() error

	//Stop the container
	stop() error

	//Kill the container
	kill() error

	//Remove the container, deleting its resources
	remove() error

	//Excecute a command within the container
	exec(command string) (string, error)

	//Get the current status of the container
	status() string

	//Get container ID
	id() string

	//Create a new container instance
	create(config ContainerConfig) (Container, error)
}
