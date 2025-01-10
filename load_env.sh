#!/bin/bash

while IFS= read -r line; do
    # Skip empty lines and comments
    [[ $line =~ ^[[:space:]]*$ ]] && continue
    [[ $line =~ ^[[:space:]]*# ]] && continue

    # Extract variable name and value
    if [[ $line =~ ^([^=]+)=(.*)$ ]]; then
        name=${BASH_REMATCH[1]}
        value=${BASH_REMATCH[2]}
        
        # Trim whitespace
        name=$(echo "$name" | xargs)
        value=$(echo "$value" | xargs)
        
        # Export variable
        export "$name=$value"
        echo "Set $name"
    fi
done < .env