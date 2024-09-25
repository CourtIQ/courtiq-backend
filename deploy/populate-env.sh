#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Enable debugging (optional). Uncomment the next line to enable debug mode.
# set -x

# Ensure the script is executed from its own directory
cd "$(dirname "$0")"

# Define the input and output files
ENV_FILE=".env"
POPULATED_ENV_FILE=".env.populated"
TARGET_DIR="../api-gateway"
SERVICE_ACCOUNT_PATH="$TARGET_DIR/firebase-service-account.json"

# Check if the 1Password CLI (op) is installed
if ! command -v op &> /dev/null
then
    echo "Error: 1Password CLI (op) is not installed. Please install it from https://developer.1password.com/docs/cli/get-started#install"
    exit 1
fi

# Check if the user is signed in to 1Password CLI
if ! op account list &> /dev/null; then
    echo "Signing in to 1Password CLI..."
    op signin
fi

# Check if the .env file exists
if [ ! -f "$ENV_FILE" ]; then
    echo "Error: $ENV_FILE file not found in $(pwd)"
    exit 1
fi

# Create or overwrite the .env.populated file
> "$POPULATED_ENV_FILE"

echo "Populating environment variables from $ENV_FILE..."

# Read the .env file line by line
while IFS='=' read -r key value || [ -n "$key" ]
do
    # Trim whitespace
    key=$(echo "$key" | xargs)
    value=$(echo "$value" | xargs)

    # Skip empty lines or lines starting with #
    if [[ -z "$key" || "$key" == \#* ]]; then
        continue
    fi

    # Check if the value starts with 'op://'
    if [[ "$value" == op://* ]]; then
        # Extract vault, item, and field from the op:// URL
        op_path=${value#op://}
        vault=$(echo "$op_path" | cut -d'/' -f1)
        item=$(echo "$op_path" | cut -d'/' -f2)
        field=$(echo "$op_path" | cut -d'/' -f3-)

        if [[ -z "$vault" || -z "$item" || -z "$field" ]]; then
            echo "Error: Invalid op:// reference '$value' for key '$key'. Expected format: op://vault/item/field"
            exit 1
        fi

        # Fetch the secret using op read
        secret=$(op read "$value" 2>/dev/null || true)

        # Check if the secret was successfully fetched
        if [ -z "$secret" ]; then
            echo "Error: Could not fetch secret for $key from 1Password. Ensure the path '$value' exists and is accessible."
            exit 1
        fi

        # Special handling for FIREBASE_ADMIN_SECRET
        if [ "$key" == "FIREBASE_ADMIN_SECRET" ]; then
            # Validate JSON
            if ! echo "$secret" | jq empty > /dev/null 2>&1; then
                echo "Error: FIREBASE_ADMIN_SECRET is not valid JSON"
                exit 1
            fi

            # Save the secret to the api-gateway directory
            echo "$secret" > "$SERVICE_ACCOUNT_PATH"
            echo "Firebase service account JSON saved to $SERVICE_ACCOUNT_PATH"
        else
            # For other secrets, handle as before
            escaped_secret=$(echo "$secret" | sed 's/\\/\\\\/g; s/"/\\"/g')
            escaped_secret=$(echo "$escaped_secret" | awk '{printf "%s\\n", $0}')
            escaped_secret=${escaped_secret%\\n}
            echo "$key=\"$escaped_secret\"" >> "$POPULATED_ENV_FILE"
        fi
    else
        # For non-1Password variables, append them as-is
        echo "$key=$value" >> "$POPULATED_ENV_FILE"
    fi
done < "$ENV_FILE"

# Construct MongoDB URI only if username and password are set
MONGO_USERNAME=$(grep '^MONGO_USERNAME=' "$POPULATED_ENV_FILE" | cut -d '=' -f2- | tr -d '"')
MONGO_PASSWORD=$(grep '^MONGO_PASSWORD=' "$POPULATED_ENV_FILE" | cut -d '=' -f2- | tr -d '"')
MONGO_HOST=$(grep '^MONGO_HOST=' "$POPULATED_ENV_FILE" | cut -d '=' -f2-)
MONGO_OPTIONS=$(grep '^MONGO_OPTIONS=' "$POPULATED_ENV_FILE" | cut -d '=' -f2-)

if [[ -n "$MONGO_USERNAME" && -n "$MONGO_PASSWORD" ]]; then
    MONGO_URI="mongodb+srv://$MONGO_USERNAME:$MONGO_PASSWORD@$MONGO_HOST/?$MONGO_OPTIONS"
    echo "MONGO_URI=\"$MONGO_URI\"" >> "$POPULATED_ENV_FILE"
else
    echo "Warning: MONGO_USERNAME or MONGO_PASSWORD is empty. MONGO_URI may be invalid."
    echo "MONGO_URI=\"mongodb+srv://:$MONGO_HOST/?$MONGO_OPTIONS\"" >> "$POPULATED_ENV_FILE"
fi

echo ".env.populated file has been successfully created."
