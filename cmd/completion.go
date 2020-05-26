/*
Copyright © 2020 mohit-kumar-sharma <flashtaken1@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/lithammer/dedent"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	completionShells = map[string]func(out io.Writer, cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

// completionCmd represents the completion command.
var (
	completionLong = dedent.Dedent(`
		Output shell completion code for the specified shell (bash/zsh).
		The shell code must be evaluated to provide interactive
		completion of proffer commands. This can be done by sourcing it from
		the .bash_profile.

		Note: this requires the bash-completion framework.

		To install it on Mac use homebrew:
			$ brew install bash-completion

		Once installed, bash_completion must be evaluated. This can be done by adding the
		following line to the .bash_profile
			$ source $(brew --prefix)/etc/bash_completion

		If bash-completion is not installed on Linux, please install the 'bash-completion' package
		via your distribution's package manager.

		Note for zsh users: [1] zsh completions are only supported in versions of zsh >= 5.2`)

	completionExample = dedent.Dedent(`
		For bash-completion
		#1 Install bash completion on a Mac using homebrew
			$ brew install bash-completion
			$ printf "\n# Bash completion support\nsource $(brew --prefix)/etc/bash_completion\n" >> $HOME/.bash_profile
			$ source $HOME/.bash_profile

		#2 Load the proffer completion code for bash into the current shell
			$ source <(proffer completion bash)

		#3 Write bash completion code to a file and source it from .bash_profile
			$ proffer completion bash > ~/.proffer/proffer_completion.bash.inc
			$ printf "\n# Proffer shell completion\nsource '$HOME/.proffer/proffer_completion.bash.inc'\n" >> $HOME/.bash_profile
			$ source $HOME/.bash_profile
		
		For zsh-completion
		# Load the proffer completion code for zsh[1] into the current shell
			$ source <(proffer completion zsh)`)

	zshInitialization = `
__proffer_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__proffer_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__proffer_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__proffer_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__proffer_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__proffer_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__proffer_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__proffer_filedir() {
	local RET OLD_IFS w qw
	__proffer_debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __proffer_debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__proffer_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__proffer_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__proffer_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__proffer_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__proffer_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__proffer_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__proffer_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__proffer_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__proffer_type/g" \
	<<'BASH_COMPLETION_EOF'
`

	completionCmd = &cobra.Command{
		Use:     "completion SHELL",
		Short:   "Generates shell completion script for specified shell type",
		Long:    completionLong,
		Example: completionExample,
		Args: func(cmd *cobra.Command, args []string) error {
			if l := len(args); l < 1 {
				return errors.New("shell type not specified")
			} else if l > 1 {
				return errors.New("too many arguments. expected only the shell type")
			}

			return nil
		},
		ValidArgs: GetSupportedShells(),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				return err
			}

			path := home + "/.proffer"
			_, err = os.Stat(path)
			if os.IsNotExist(err) {
				if err := os.Mkdir(path, 0755); err != nil {
					return err
				}
			}

			return RunCompletion(os.Stdout, cmd, args)
		},
	}
)

func init() {
	rootCmd.AddCommand(completionCmd)
}

// GetSupportedShells returns a list of supported shells.
func GetSupportedShells() []string {
	shells := []string{}
	for s := range completionShells {
		shells = append(shells, s)
	}

	return shells
}

func runCompletionBash(out io.Writer, proffer *cobra.Command) error {
	return proffer.GenBashCompletion(out)
}

func runCompletionZsh(out io.Writer, proffer *cobra.Command) error {
	_, err := out.Write([]byte(zshInitialization))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := proffer.GenBashCompletion(buf); err != nil {
		return err
	}

	_, err = out.Write(buf.Bytes())
	if err != nil {
		return err
	}

	zshTail := `
BASH_COMPLETION_EOF
}
__proffer_bash_source <(__proffer_convert_bash_to_zsh)
`
	_, err = out.Write([]byte(zshTail))

	return err
}

// RunCompletion checks given arguments and executes command.
func RunCompletion(out io.Writer, cmd *cobra.Command, args []string) error {
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("unsupported shell type %q", args[0])
	}

	return run(out, cmd.Parent())
}
