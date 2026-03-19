package completion

// generatePowerShell returns the PowerShell completion script.
func generatePowerShell() string {
	return `Register-ArgumentCompleter -CommandName gitmap -ScriptBlock {
    param($wordToComplete, $commandAst, $cursorPosition)
    $elems = $commandAst.CommandElements | Select-Object -Skip 1
    $cmd = if ($elems.Count -gt 0) { $elems[0].ToString() } else { "" }
    $prev = if ($elems.Count -gt 1) { $elems[$elems.Count - 1].ToString() } else { "" }

    if ($cmd -eq "cd" -or $cmd -eq "go") {
        if ($prev -eq "--group" -or $prev -eq "-g") {
            $items = gitmap completion --list-groups
        } else {
            $items = @(gitmap completion --list-repos) + @("repos", "set-default", "clear-default")
        }
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "pull") {
        gitmap completion --list-repos | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "exec" -and ($prev -eq "--group")) {
        gitmap completion --list-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "group" -or $cmd -eq "g") {
        $subs = @("create", "add", "remove", "list", "show", "delete", "pull", "status", "exec", "clear")
        $groups = @(gitmap completion --list-groups)
        $items = $subs + $groups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "list" -or $cmd -eq "ls") {
        $items = @("go", "node", "nodejs", "react", "cpp", "csharp", "groups", "--group", "--verbose")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "multi-group" -or $cmd -eq "mg") {
        $subs = @("pull", "status", "exec", "clear")
        $groups = @(gitmap completion --list-groups)
        $items = $subs + $groups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "release" -or $cmd -eq "r") {
        $items = @("--assets", "--commit", "--branch", "--bump", "--draft", "--dry-run", "--compress", "--checksums", "--no-assets", "--targets", "--list-targets", "--verbose", "--zip-group", "-Z", "--bundle")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "release-branch" -or $cmd -eq "rb") {
        $items = @("--assets", "--draft", "--dry-run", "--compress", "--checksums", "--no-assets", "--targets")
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "alias" -or $cmd -eq "a") {
        $subs = @("set", "remove", "list", "show", "suggest")
        $aliases = @(gitmap completion --list-aliases)
        $items = $subs + $aliases
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($cmd -eq "zip-group" -or $cmd -eq "z") {
        $subs = @("create", "add", "remove", "list", "show", "delete", "rename")
        $zgroups = @(gitmap completion --list-zip-groups)
        $items = $subs + $zgroups
        $items | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($prev -eq "-A" -or $prev -eq "--alias") {
        gitmap completion --list-aliases | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    if ($prev -eq "--zip-group") {
        gitmap completion --list-zip-groups | Where-Object { $_ -like "$wordToComplete*" } |
            ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
        return
    }

    gitmap completion --list-commands | Where-Object { $_ -like "$wordToComplete*" } |
        ForEach-Object { [System.Management.Automation.CompletionResult]::new($_) }
}
`
}
