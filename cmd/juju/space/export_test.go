// Copyright 2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package space

import "github.com/juju/juju/cmd/juju/subnet"

func NewCreateCommand(api SpaceAPI) *CreateCommand {
	createCmd := &CreateCommand{}
	createCmd.api = api
	return createCmd
}

func NewRemoveCommand(api SpaceAPI) *RemoveCommand {
	removeCmd := &RemoveCommand{}
	removeCmd.api = api
	return removeCmd
}

func NewUpdateCommand(api SpaceAPI) *UpdateCommand {
	updateCmd := &UpdateCommand{}
	updateCmd.api = api
	return updateCmd
}

func NewRenameCommand(api SpaceAPI) *RenameCommand {
	renameCmd := &RenameCommand{}
	renameCmd.api = api
	return renameCmd
}

func NewListCommand(api SpaceAPI, subnetapi subnet.SubnetAPI) *ListCommand {
	listCmd := &ListCommand{}
	listCmd.api = api
	listCmd.subnetapi = subnetapi
	return listCmd
}

func ListFormat(cmd *ListCommand) string {
	return cmd.out.Name()
}
