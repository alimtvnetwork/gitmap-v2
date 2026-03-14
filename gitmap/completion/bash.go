package completion

// generateBash returns the Bash completion script.
func generateBash() string {
	return `_gitmap_completions() {
    local cur prev cmd
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    cmd="${COMP_WORDS[1]}"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "$(gitmap completion --list-commands)" -- "$cur"))
        return
    fi

    case "$cmd" in
        cd|go)
            if [[ "$prev" == "--group" || "$prev" == "-g" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-groups)" -- "$cur"))
            else
                COMPREPLY=($(compgen -W "$(gitmap completion --list-repos) repos set-default clear-default" -- "$cur"))
            fi
            ;;
        pull)
            COMPREPLY=($(compgen -W "$(gitmap completion --list-repos)" -- "$cur"))
            ;;
        exec)
            if [[ "$prev" == "--group" ]]; then
                COMPREPLY=($(compgen -W "$(gitmap completion --list-groups)" -- "$cur"))
            fi
            ;;
        group|g)
            COMPREPLY=($(compgen -W "create add remove list show delete pull status exec clear $(gitmap completion --list-groups)" -- "$cur"))
            ;;
        list|ls)
            COMPREPLY=($(compgen -W "go node nodejs react cpp csharp groups --group --verbose" -- "$cur"))
            ;;
        multi-group|mg)
            COMPREPLY=($(compgen -W "pull status exec clear $(gitmap completion --list-groups)" -- "$cur"))
            ;;
    esac
}
complete -F _gitmap_completions gitmap
`
}
