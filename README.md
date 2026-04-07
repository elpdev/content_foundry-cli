# Content Foundry CLI

A command-line interface for the [Content Foundry](https://github.com/elpdev/content_foundry) API. Manage your entire content pipeline from the terminal -- sources, content items, drafts, AI personas, media generation, and more.

Built with [Cobra](https://github.com/spf13/cobra) and the [Charm](https://charm.sh) stack (lipgloss, huh).

## Install

### Homebrew

```sh
brew tap elpdev/tap
brew install content_foundry
```

### From source

Requires Go 1.22+.

```sh
git clone soft-serve:content_foundry-cli.git
cd content_foundry-cli
make install
```

## Quick Start

```sh
# 1. Connect to your Content Foundry server
content_foundry auth login

# 2. List your brands
content_foundry brands list

# 3. Set a default brand
content_foundry brands use 1

# 4. You're ready -- list sources, drafts, etc.
content_foundry sources list
content_foundry drafts list --status pending_review
```

## Authentication

The CLI uses JWT authentication with API keys. Run `content_foundry auth login` to enter your credentials interactively:

- **Base URL** -- your Content Foundry server (e.g. `http://localhost:3000`)
- **Client ID** -- your API key client_id
- **Secret Key** -- your API key secret

Credentials are saved to `~/.config/content_foundry/config.toml`. JWT tokens are cached in `~/.config/content_foundry/token.toml` and auto-refresh on expiry.

## Commands

### Content Pipeline

```sh
# Sources -- ingest content from RSS, Twitter, etc.
content_foundry sources list [--active true]
content_foundry sources show <id>
content_foundry sources create --name "Blog RSS" --type "Sources::RssFeed" --config '{"feed_url":"https://..."}'
content_foundry sources update <id> --name "..." [--active false]
content_foundry sources delete <id>
content_foundry sources fetch <id>                  # trigger async fetch

# Content Items -- ingested content
content_foundry content-items list [--status pending]
content_foundry content-items show <id>             # includes associated drafts
content_foundry content-items process <id> [--guidance "..."]
content_foundry content-items generate-drafts <id>  # auto-generate for all platforms
```

### Drafts Workflow

```sh
# CRUD
content_foundry drafts list [--status pending_review] [--platform-id 1] [--assigned-to 5]
content_foundry drafts show <id>                    # includes comments + publication
content_foundry drafts create --title "..." --content "..." --platform-id 1
content_foundry drafts delete <id>

# Review workflow
content_foundry drafts approve <id>
content_foundry drafts reject <id>
content_foundry drafts revise <id> --notes "Please shorten the intro"

# Scheduling
content_foundry drafts schedule <id> --at "2026-04-10T14:00:00Z"
content_foundry drafts reschedule <id> --at "2026-04-11T10:00:00Z"
content_foundry drafts unschedule <id>

# Assignment
content_foundry drafts assign <id> --to <user_id>
content_foundry drafts unassign <id>

# Media + Publishing
content_foundry drafts media <id> --turn-ids 10,11,12
content_foundry drafts publish <id> [--turn-ids 10,11]

# Comments
content_foundry drafts comments add <draft_id> --body "Looks good!"
content_foundry drafts comments delete <draft_id> <comment_id>
```

### AI (Personas & Chats)

```sh
# Personas
content_foundry personas list
content_foundry personas show <id>
content_foundry personas create --name "Social Media Manager" --role "Content Creator" --prompt "..."
content_foundry personas update <id> --name "..."
content_foundry personas delete <id>

# Chats
content_foundry chats list [--persona-id 1]
content_foundry chats show <id>                     # shows full message history
content_foundry chats create --prompt "Write a thread about..." [--persona-id 1] [--model "..."]
content_foundry chats delete <id>
content_foundry chats send <id> --message "Can you make it more casual?"
```

### Media Generation

```sh
# Media sessions
content_foundry media list [--type image]
content_foundry media show <id>                     # lists all turns
content_foundry media create --type image --prompt "A sunset over mountains" [--aspect-ratio "16:9"]
content_foundry media create --type video --prompt "Ocean waves" [--duration 8]
content_foundry media create --type audio --prompt "Lo-fi beat" [--seconds 30]
content_foundry media delete <id>

# Image turns
content_foundry images show <session_id> <turn_id>  # renders inline if terminal supports it
content_foundry images create <session_id> --prompt "..."
content_foundry images convert <session_id> <turn_id>
content_foundry images delete <session_id> <turn_id>

# Video turns
content_foundry videos show <session_id> <turn_id>
content_foundry videos create <session_id> --prompt "..." [--duration 8]
content_foundry videos extend <session_id> <turn_id> --target-duration 30
content_foundry videos delete <session_id> <turn_id>

# Audio turns
content_foundry audio create <session_id> --prompt "..." [--seconds 30]
content_foundry audio delete <session_id> <turn_id>
```

### Supporting Resources

```sh
# Brands
content_foundry brands list
content_foundry brands show <id>
content_foundry brands create --name "My Brand" [--slug "my-brand"]
content_foundry brands update <id> --name "..." [--voice "..."] [--mission "..."] [--target-audience "..."]
content_foundry brands delete <id>
content_foundry brands use <id|slug>                # set default brand

# Platforms
content_foundry platforms list [--active true]
content_foundry platforms show <id>
content_foundry platforms create --name "Twitter" --type "Platforms::Twitter"
content_foundry platforms update <id> [--name "..."] [--active false]
content_foundry platforms delete <id>

# Labels (per platform)
content_foundry platforms labels list <platform_id>
content_foundry platforms labels show <platform_id> <label_id>
content_foundry platforms labels create <platform_id> --name "Politics" [--slug "politics"]
content_foundry platforms labels update <platform_id> <label_id> --name "..."
content_foundry platforms labels delete <platform_id> <label_id>

# Assets
content_foundry assets list [--type image]
content_foundry assets upload <file_path>
content_foundry assets delete <id>

# Brand Documents
content_foundry docs list
content_foundry docs show <id>
content_foundry docs create --url "https://..."
content_foundry docs delete <id>
content_foundry docs index-content

# Activity
content_foundry activity list

# Notifications
content_foundry notifications list
content_foundry notifications read <id>
content_foundry notifications read-all

# Tickets
content_foundry tickets list [--category bug_report] [--priority high] [--unresolved]
content_foundry tickets show <id>
content_foundry tickets create --subject "..." --description "..." --category bug_report --priority high
content_foundry tickets update <id> [--priority medium]

# Account
content_foundry account show
content_foundry account update --name "..."

# Members
content_foundry members list
content_foundry members show <id>
content_foundry members update <id> --admin --editor
content_foundry members delete <id>
content_foundry members brand-access <id> --brand-ids 1,2,3

# Invitations
content_foundry invitations list
content_foundry invitations show <token>
content_foundry invitations create --name "Jane" --email "jane@co.com" [--admin] [--brand-ids 1,2]
content_foundry invitations delete <token>
```

## Global Flags

| Flag                 | Description                                                         |
| -------------------- | ------------------------------------------------------------------- |
| `--config <path>`    | Config file path (default: `~/.config/content_foundry/config.toml`) |
| `-f, --format`       | Output format: `table` (default), `json`, `text`                    |
| `-v, --verbose`      | Enable verbose logging                                              |
| `--brand <id\|slug>` | Override the default brand for this command                         |

## Output Formats

```sh
# Styled table (default)
content_foundry brands list

# JSON (useful for piping to jq)
content_foundry brands list -f json

# Plain text (useful for scripting)
content_foundry brands list -f text
```

## Inline Image Rendering

The CLI can render images directly in your terminal when viewing image turns. Supported terminals:

- **iTerm2** (macOS)
- **WezTerm**
- **Kitty**

Other terminals will display image metadata without inline rendering.

## Project Structure

```
content_foundry_cli/
  main.go                      # Entry point
  Makefile                     # build, install, clean, lint, test
  cmd/                         # Cobra commands (one file per resource)
  internal/
    api/                       # HTTP client + service layer (one file per resource)
    config/                    # TOML config + JWT token caching
    models/                    # Go structs matching API JSON responses
    output/                    # Formatter interface (table, json, text)
    media/                     # Terminal image rendering (iTerm2, Kitty protocols)
    logging/                   # Structured logging (charmbracelet/log)
```

## Development

```sh
make build     # compile binary
make lint      # go vet
make test      # go test
make run       # go run
```

## Releasing a New Version

Merges to `main` automatically create the next patch release, publish GitHub release artifacts, and update the Homebrew tap formula.

```sh
# Upgrade locally after a release lands
brew upgrade content_foundry
```

Repository setup requirements:

- Add the `HOMEBREW_TAP_TOKEN` Actions secret with write access to `elpdev/homebrew-tap`
- Keep the tap formula in `elpdev/homebrew-tap/Formula/content_foundry.rb`

## License

Private.
