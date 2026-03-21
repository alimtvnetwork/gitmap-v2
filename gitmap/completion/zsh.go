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
            if [[ "${words[CURRENT-1]}" == "--zip-group" ]]; then
                local -a zgroups=($(gitmap completion --list-zip-groups))
                _describe 'zip-group' zgroups
            else
                local -a flags=(--assets --commit --branch --bump --draft --dry-run --compress --checksums --no-assets --targets --list-targets --verbose --zip-group -Z --bundle)
                _describe 'flag' flags
            fi
            ;;
        release-branch|rb)
            local -a flags=(--assets --draft --dry-run --compress --checksums --no-assets --targets)
            _describe 'flag' flags
            ;;
        alias|a)
            local -a subs=(set remove list show suggest)
            local -a aliases=($(gitmap completion --list-aliases))
            subs+=("${aliases[@]}")
            _describe 'subcommand' subs
            ;;
        zip-group|z)
            if (( CURRENT >= 4 )) && [[ "${words[3]}" == "add" || "${words[3]}" == "show" || "${words[3]}" == "delete" || "${words[3]}" == "remove" || "${words[3]}" == "rename" ]]; then
                local -a zgroups=($(gitmap completion --list-zip-groups))
                _describe 'zip-group' zgroups
            else
                local -a subs=(create add remove list show delete rename)
                local -a zgroups=($(gitmap completion --list-zip-groups))
                subs+=("${zgroups[@]}")
                _describe 'subcommand' subs
            fi
            ;;
        *)
            if [[ "${words[CURRENT-1]}" == "-A" || "${words[CURRENT-1]}" == "--alias" ]]; then
                local -a aliases=($(gitmap completion --list-aliases))
                _describe 'alias' aliases
            elif [[ "${words[CURRENT-1]}" == "--zip-group" ]]; then
                local -a zgroups=($(gitmap completion --list-zip-groups))
                _describe 'zip-group' zgroups
            fi
            ;;
    esac
}

compdef _gitmap gitmap
`
}
