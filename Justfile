_default:
    @just --list

# Link the AI directories
ai-link:
    @mkdir -p .{cursor,claude}
    @stow --dir=ai --target=.cursor .
    @stow --dir=ai --target=.claude .