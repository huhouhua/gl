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

package file

import (
	"fmt"
	"github.com/AlekSi/pointer"
	cmdutil "github.com/huhouhua/gitlab-repo-operator/cmd/util"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"
	"strings"
)

type ListOptions struct {
	gitlabClient *gitlab.Client
	file         *gitlab.ListTreeOptions
	project      string
	path         string
	Out          string
	All          bool
}

func NewListOptions() *ListOptions {
	return &ListOptions{
		file: &gitlab.ListTreeOptions{
			ListOptions: gitlab.ListOptions{
				Page:    1,
				PerPage: 50,
			},
			Path:      pointer.ToString(""),
			Recursive: pointer.ToBool(true),
		},
		Out: "simple",
	}
}

var (
	getFilesDesc = "get file for project "

	getFilesExample = `# list project file
grepo get files myProject`
)

func NewGetFilesCmd(f cmdutil.Factory) *cobra.Command {
	o := NewListOptions()
	cmd := &cobra.Command{
		Use:                   "files",
		Aliases:               []string{"f"},
		Short:                 getFilesDesc,
		Example:               getFilesExample,
		DisableFlagsInUseLine: true,
		TraverseChildren:      true,
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(f, cmd, args))
			cmdutil.CheckErr(o.Validate(cmd, args))
			cmdutil.CheckErr(o.Run(args))
		},
		SuggestFor: []string{"file"},
	}
	return cmd
}

func (o *ListOptions) AddFlags(cmd *cobra.Command) {
	cmdutil.AddPaginationVarFlags(cmd, &o.file.ListOptions)
	cmdutil.AddOutFlag(cmd, &o.Out)
	cmdutil.AddSortVarFlag(cmd, &o.file.Sort)
	f := cmd.Flags()
	f.StringVar(o.file.Ref, "ref", "", "The name of a repository branch or tag or, if not given, the default branch.")
	f.StringVarP(o.file.Path, "path", "p", *o.file.Path, "The path inside the repository. Used to get content of subdirectories.")
	f.BoolVarP(o.file.Recursive, "recursive", "r", *o.file.Recursive, "Boolean value used to get a recursive tree. Default is true.")
	f.BoolVarP(&o.All, "all", "A", o.All, "If present, list the across all project file. file in current context is ignored even if specified with --all.")
}

// Complete completes all the required options.
func (o *ListOptions) Complete(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	if len(args) > 0 {
		o.project = args[0]
	}
	if o.file.Ref != nil && strings.TrimSpace(*o.file.Ref) == "" {
		o.file.Ref = nil
	}
	o.gitlabClient, err = f.GitlabClient()
	return err
}

// Validate makes sure there is no discrepency in command options.
func (o *ListOptions) Validate(cmd *cobra.Command, args []string) error {
	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return fmt.Errorf("please enter project name and id")
	}
	return nil
}

// Run executes a list subcommand using the specified options.
func (o *ListOptions) Run(args []string) error {
	var list []*gitlab.TreeNode
	if o.All {
		o.file.ListOptions.PerPage = 100
		o.file.ListOptions.Page = 1
	}
	for {
		tree, _, err := o.gitlabClient.Repositories.ListTree(o.project, o.file)
		if err != nil {
			return nil
		}
		list = append(list, tree...)
		if cap(tree) == 0 || !o.All {
			break
		}
		o.file.ListOptions.Page++
	}
	cmdutil.PrintFilesOut(o.Out, list...)
	return nil
}
