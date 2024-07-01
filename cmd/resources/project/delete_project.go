// Copyright 2024 The huhouhua Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http:www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package project

import (
	"fmt"
	"github.com/huhouhua/gl/cmd/require"
	cmdutil "github.com/huhouhua/gl/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type DeleteOptions struct {
	gitlabClient *gitlab.Client
	project      string
	ioStreams    cmdutil.IOStreams
}

var (
	deleteProjectDesc = "Delete a Gitlab project by specifying the full path"

	deleteProjectExample = `# delete a project
gl delete project ProjectX

# delete a project under a group
gl delete project group/project`
)

func NewDeleteOptions(ioStreams cmdutil.IOStreams) *DeleteOptions {
	return &DeleteOptions{
		ioStreams: ioStreams,
	}
}

func NewDeleteProjectCmd(f cmdutil.Factory, ioStreams cmdutil.IOStreams) *cobra.Command {
	o := NewDeleteOptions(ioStreams)
	cmd := &cobra.Command{
		Use:                   "project",
		Aliases:               []string{"p"},
		Short:                 deleteProjectDesc,
		Example:               deleteProjectExample,
		Args:                  require.ExactArgs(1),
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{},
	}
	return cmd
}

// Complete completes all the required options.
func (o *DeleteOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	o.gitlabClient, err = f.GitlabClient()
	if len(args) > 0 {
		o.project = args[0]
	}
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *DeleteOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("please enter projectName and id")
	}
	if strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("error from server (NotFound): project %s not found", args[0])
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *DeleteOptions) Run(args []string) error {
	projectInfo, _, err := o.gitlabClient.Projects.GetProject(o.project, &gitlab.GetProjectOptions{})
	if err != nil {
		return err
	}
	_, err = o.gitlabClient.Projects.DeleteProject(projectInfo.ID)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(o.ioStreams.Out, "project (%s) with id (%d) has been deleted\n", o.project, projectInfo.ID)
	return nil
}
