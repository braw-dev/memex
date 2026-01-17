_default:
    @just --list

# Link the AI directories
ai-link:
    @mkdir -p .{cursor,claude}
    @stow --dir=ai --target=.cursor .
    @stow --dir=ai --target=.claude .

# Install the dev dependencies
install-dev:
    @go mod tidy
    @go get -tool github.com/evilmartians/lefthook/v2
    @lefthook install