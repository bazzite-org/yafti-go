default:
    @just --list

dev:
    go tool templ generate --watch --cmd="env YAFTI_CONF=\"$PWD/example_config.yml\" go run ."
