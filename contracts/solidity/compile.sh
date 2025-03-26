#!/bin/bash

# Exit on any error
set -e

# Set default flag values
no_deps=false
input_file=""

# Parse arguments
for arg in "$@"; do
    if [ "$arg" = "--no-deps" ]; then
        no_deps=true
    elif [[ "$arg" =~ \.sol$ ]]; then
        input_file="$arg"
    fi
done

# Check if an input file was provided
if [ -z "$input_file" ]; then
    echo "Error: Please provide a Solidity file name."
    exit 1
fi

# Extract the filename without extension
input_file_name="${input_file##*/}"
input_file_name="${input_file_name%.*}"

# Extract the directory of the input file
input_dir="$(dirname "$input_file")"

# Set the output directory
output_dir=".temp"

# Create the output directory if it doesn't exist
mkdir -p "$output_dir"

# Compile the Solidity file and output to the specified directory
# Include node_modules in the include path
solcjs "$input_file" --bin --abi --optimize --base-path . --include-path node_modules -o "$output_dir"

# Rename the files
for file in "$output_dir"/*_sol_*.abi "$output_dir"/*_sol_*.bin; do
    # Extract the filename without the prefix and suffix
    filename=${file##*/}
    filename=${filename#*_sol_}
    filename=${filename%.*}

    # Determine the extension
    extension=${file##*.}

    # Rename the file
    mv "$file" "$output_dir/$filename.$extension"
done

# Move the renamed files to the input directory
if $no_deps; then
    # Only move files for the target contract
    mv "$output_dir/$input_file_name.abi" "$output_dir/$input_file_name.bin" "$input_dir"
else
    # Move all files
    mv "$output_dir"/*.abi "$output_dir"/*.bin "$input_dir"
fi

# Delete the temporary output directory
rm -rf "$output_dir"
