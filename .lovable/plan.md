## Spec: CD Page, Alias Improvements & CodeBlock Font Size

### 1. New CD (go) Documentation Page — `/cd`

**Purpose**: Document `gitmap cd` (alias `go`) — navigate to repo directories by name with fuzzy matching.

**Sections**:
- **Terminal Preview**: Shows `gitmap cd myrepo` jumping to directory, and a failed match showing fuzzy suggestions
- **Features**: Direct navigation, fuzzy suggestion on mismatch, interactive picker (`repos`), group filtering, shell wrapper (`gcd`)
- **Usage**: `gitmap cd <repo-name|repos> [--group <name>] [--pick]`
- **Flags**: `--group`, `--pick`, `--verbose`
- **Fuzzy Matching**: When exact name not found → shows "Did you mean?" with closest matches ranked by Levenshtein distance
- **Shell Wrapper**: `gcd` function auto-installed by `gitmap setup`
- **Examples**: Direct jump, picker, group-scoped, alias integration (`-A`)
- **See Also**: list, scan, group, alias

### 2. Alias Page Improvements

**Current state**: Has basic structure but missing help-style examples.

**Changes**:
- Add `--help` output as a CodeBlock (the actual help text like helptext/alias.md would show)
- Add more realistic terminal examples with descriptions
- Apply `docs-h2`, `docs-table` CSS classes consistently
- Add fuzzy matching behavior when alias not found

### 3. CodeBlock Font Size Controls

**Changes**:
- Increase default code font from `text-sm` to `text-[15px]`
- Add font size toggle button (A↑/A↓) in the CodeBlock header
- Three sizes: small (13px), medium (15px), large (17px)
- Persist choice in component state

### Files to create/modify:
- `src/pages/Cd.tsx` — NEW
- `src/pages/Alias.tsx` — UPDATE
- `src/components/docs/CodeBlock.tsx` — UPDATE (font size)
- `src/App.tsx` — Add route for `/cd`
- `src/components/docs/DocsSidebar.tsx` — Add CD nav item
