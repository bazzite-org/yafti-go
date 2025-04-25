default:
    @just --list

dev:
    go tool templ generate --watch --cmd='go run .'
