#!/bin/bash

input_folder="./input"

# Check if the input folder exists
if [ ! -d "$input_folder" ]; then
    echo "Error: Input folder does not exist!"
    exit 1
fi

# Iterate over each file in the input folder
for file in "$input_folder"/*; do
    # Check if it is a regular file
    if [ -f "$file" ]; then
        # Write "true" to the file
        echo "" > "$file"
    fi
done
