/*
This file is part of the Gadget command-line tools.
Copyright (C) 2017 Next Thing Co.

Gadget is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 2 of the License, or
(at your option) any later version.

Gadget is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with Gadget.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"fmt"
	"github.com/nextthingco/libgadget"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func addUsage() error {
	log.Info("Usage: gadget [flags] add [type] [name]")
	log.Info("               *opt        *req   *req ")
	log.Info("Type: service | onboot                 ")
	log.Info("Name: friendly name for container      ")

	return errors.New("Incorrect add usage")
}

// Process the build arguments and execute build
func GadgetAdd(args []string, g *libgadget.GadgetContext) error {

	addUu := uuid.NewV4()

	if len(args) != 2 {
		return addUsage()
	}

	log.Infof("Adding new %s: \"%s\" ", args[0], args[1])

	addGadgetContainer := libgadget.GadgetContainer{
		Name:  args[1],
		Image: fmt.Sprintf("%s/%s", g.Config.Name, args[1]),
		UUID:  fmt.Sprintf("%s", addUu),
	}

	// parse arguments
	switch args[0] {
	case "service":
		g.Config.Services = append(g.Config.Services, addGadgetContainer)
	case "onboot":
		g.Config.Onboot = append(g.Config.Onboot, addGadgetContainer)
	default:
		log.Errorf("  %q is not valid command.", args[0])
		return addUsage()
	}

	g.Config = libgadget.CleanConfig(g.Config)

	fileLocation := fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory)

	outBytes, err := yaml.Marshal(g.Config)
	if err != nil {

		log.WithFields(log.Fields{
			"function":   "GadgetAdd",
			"location":   fileLocation,
			"init-stage": "parsing",
		}).Debug("The config file is probably malformed")

		log.Errorf("Failed to parse config file [%s]", fileLocation)
		log.Warn("Is this a valid gadget.yaml?")
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/gadget.yml", g.WorkingDirectory), outBytes, 0644)
	if err != nil {

		log.WithFields(log.Fields{
			"function":   "GadgetAdd",
			"location":   fileLocation,
			"init-stage": "writing file",
		}).Debug("This is likely due to a problem with permissions")

		log.Errorf("Failed to edit config file [%s]", fileLocation)
		log.Warn("Do you have permission to modify this file?")

		return err
	}

	return err
}
