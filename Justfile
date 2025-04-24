default:
    @just --list

dev:
    #!/bin/bash
    set -emo pipefail
    jobs_pids=()

    tailwindcss --watch -i ./static/css/tailwind.css -o ./static/css/main.css &
    jobs_pids+=($!)

    go tool templ generate --watch --cmd='go run .' &
    jobs_pids+=($!)

    trap 'echo "${jobs_pids[@]}" | xargs kill' INT

    echo "Running processes in the background. Press Ctrl+C to exit."
    wait "${jobs_pids[@]}"
