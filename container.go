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
	Start() error

	//Stop the container
	Stop() error

	//Kill the container
	Kill() error

	//Remove the container, deleting its resources
	Remove() error

	//Excecute a command within the container
	Exec(command string) (string, error)

	//Get the current status of the container
	Status() string

	//Get container ID
	ID() string

	//Create a new container instance
	Create(config ContainerConfig) (Container, error)
}
