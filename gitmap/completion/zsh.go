package completion

// generateZsh returns the Zsh completion script.
func generateZsh() string {
	return `#compdef gitmap

_gitmap() {
    local -a commands repos groups

    if (( CURRENT == 2 )); then
        commands=($(gitmap completion --list-commands))
        _describe 'command' commands
        return
    fi

    case "${words[2]}" in
        cd|go)
            if [[ "${words[CURRENT-1]}" == "--group" || "${words[CURRENT-1]}" == "-g" ]]; then
                groups=($(gitmap completion --list-groups))
                _describe 'group' groups
            else
                repos=($(gitmap completion --list-repos))
                repos+=(repos set-default clear-default)
                _describe 'repo' repos
            fi
            ;;
        pull)
            repos=($(gitmap completion --list-repos))
            _describe 'repo' repos
            ;;
        exec)
            if [[ "${words[CURRENT-1]}" == "--group" ]]; then
                groups=($(gitmap completion --list-groups))
                _describe 'group' groups
            fi
            ;;
        group|g)
            local -a subs=(create add remove list show delete pull status exec clear)
            groups=($(gitmap completion --list-groups))
            subs+=("${groups[@]}")
            _describe 'subcommand' subs
            ;;
        list|ls)
            local -a types=(go node nodejs react cpp csharp groups)
            _describe 'type' types
            ;;
        multi-group|mg)
            local -a subs=(pull status exec clear)
            groups=($(gitmap completion --list-groups))
            subs+=("${groups[@]}")
            _describe 'subcommand' subs
            ;;
        release|r)
            local -a flags=(--assets --commit --branch --bump --draft --dry-run --compress --checksums --verbose)
            _describe 'flag' flags
            ;;
        release-branch|rb)
            local -a flags=(--assets --draft --dry-run --compress --checksums)
            _describe 'flag' flags
            ;;
    esac
}

compdef _gitmap gitmap
`
}
