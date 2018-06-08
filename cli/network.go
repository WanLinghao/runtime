// Copyright (c) 2017 HyperHQ Inc.
// Copyright (c) 2018 Huawei Corporation.
//
// SPDX-License-Identifier: Apache-2.0
//

package main

import (
	"fmt"

	"os"

	"encoding/json"

	"github.com/kata-containers/agent/protocols/grpc"
	vc "github.com/kata-containers/runtime/virtcontainers"
	"github.com/urfave/cli"
)

var networkCLICommand = cli.Command{
	Name:  "network",
	Usage: "manage interfaces and routes for container",
	Subcommands: []cli.Command{
		netListCommand,
		netUpdateCommand,
	},
	Action: func(context *cli.Context) error {
		return cli.ShowSubcommandHelp(context)
	},
}

var netListCommand = cli.Command{
	Name:      "ls",
	Aliases:   []string{"list"},
	Usage:     "list network interfaces and routes in a container",
	ArgsUsage: `ls <container-id>`,
	Flags:     []cli.Flag{},
	Action: func(context *cli.Context) error {
		if context.Args().Present() == false {
			return fmt.Errorf("missing container ID, should provide one")
		}
		containerID := context.Args().First()
		status, sandboxID, err := getExistingContainerInfo(containerID)
		if err != nil {
			return err
		}

		containerID = status.ID
		// container MUST be running
		if status.State.State != vc.StateRunning {
			return fmt.Errorf("container %s is not running", containerID)
		}

		return nil
	},
}

var netUpdateCommand = cli.Command{
	Name:      "update",
	Usage:     "update configuration of interface",
	ArgsUsage: `update <container-id>`,
	Flags:     []cli.Flag{},
	Action: func(context *cli.Context) error {
		if context.Args().Present() == false {
			return fmt.Errorf("missing container ID, should provide one")
		}
		containerID := context.Args().First()
		status, sandboxID, err := getExistingContainerInfo(containerID)
		if err != nil {
			return err
		}

		containerID = status.ID
		// container MUST be running
		if status.State.State != vc.StateRunning {
			return fmt.Errorf("container %s is not running", containerID)
		}

		var inf *grpc.Interface

		fileName := context.Args().Get(1)

		f, _ := os.Open(fileName)
		json.NewDecoder(f).Decode(&inf)

		vci.AddNetwork(sandboxID, inf)
		return nil
	},
}
