#! /bin/bash

app="$1"

if [ "$app" = "web" ]; then
    exec ./web "${@:2}"
elif [ "$app" = "migrate" ]; then
    exec ./migrate "${@:2}"
else
    echo "Error: Invalid application specified. Use 'web' or 'migrate'."
    exit 1
fi
