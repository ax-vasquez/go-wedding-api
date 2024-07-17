#! /bin/bash

EXPECTED_KEYS=(
    "PGSQL_HOST"
    "PGSQL_USER"
    "PGSQL_PASSWORD"
    "PGSQL_DBNAME"
    "PGSQL_PORT"
    "PGSQL_TIMEZONE"
)

value=$(<.env)

# Check if .env file exists
if ! [ -e .env ]; then
    echo "No .env file found - creating empty .env file"
    touch .env
    echo "# The contents of this file are generated when you run ./setup.sh" >> .env
    echo "# You can edit variables as-needed; they will not be overwritten when re-running the setup script" >> .env
fi

# Adds missing keys to the .env file, which is done both when updating and creating the .env file
for key in "${EXPECTED_KEYS[@]}"
do
    # If a key is missing, append it to the .env file and use default values
    if ! grep -q ${key} ".env"
    then
        if [ ${key} = "PGSQL_USER" ] || [ ${key} = "PGSQL_PASSWORD" ] || [ ${key} = "PGSQL_DBNAME" ]; then
            echo "Using \"gorm\" as default value for \"${key}\""
            echo "${key}=gorm" >> .env
        elif [ ${key} = "PGSQL_HOST" ]; then
            echo "Using \"localhost\" as default value for \"${key}\""
            echo "${key}=localhost" >> .env
        elif [ ${key} = "PGSQL_PORT" ]; then
            echo "Using \"9920\" as default value for \"${key}\""
            echo "${key}=9920" >> .env
        elif [ ${key} = "PGSQL_TIMEZONE" ]; then
            echo "Using \"US/Central\" as default value for \"${key}\""
            echo "${key}=US/Central" >> .env
        fi
    fi
done
