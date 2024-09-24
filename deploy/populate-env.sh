#!/bin/bash

# Ensure we're in the deploy directory
cd "$(dirname "$0")"

# Ensure 1Password CLI is installed and authenticated
if ! command -v op &> /dev/null
then
    echo "1Password CLI could not be found. Please install it first."
    exit 1
fi

# Check if we're already signed in
if ! op account get &> /dev/null; then
    echo "Please sign in to 1Password CLI"
    eval $(op signin)
fi

# Read the template .env file
while IFS= read -r line
do
    # Check if the line contains a command substitution
    if [[ $line == *"$(op"* ]]; then
        # Evaluate the line and add it to the new .env file
        eval echo "$line" >> .env.populated
    else
        # If not, just copy the line as is
        echo "$line" >> .env.populated
    fi
done < .env

echo ".env file has been populated with secrets from 1Password."

# Return to the original directory
cd -