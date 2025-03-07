#!/bin/bash

# Check if a directory path is passed as an argument
if [ -z "$1" ]; then
    # If no argument is provided, use the current directory
    directory="./"
else
    directory="$1"
fi

# Get the total number of SVG files in the directory
total_files=$(ls "$directory"/*.svg 2>/dev/null | wc -l)

# Exit if no SVG files are found
if [ "$total_files" -eq 0 ]; then
    echo "No SVG files found in the directory."
    exit 1
fi

# Initialize the JSON array using jq
output="[]"

# Start the progress bar loop
counter=0
for file in "$directory"/*.svg; do
    if [[ -f "$file" ]]; then
        # Extract the file name without the path
        filename=$(basename "$file")
        
        # Extract fill and stroke colors from the SVG using grep and sed
        colors=$(grep -o 'fill="[^"]*"\|stroke="[^"]*"' "$file" | sed -E 's/^(fill|stroke)="([^"]+)"/\2/' | tr '[:upper:]' '[:lower:]')
        
        # Remove any color that is "none", "white", or similar white/near-white colors
        colors=$(echo "$colors" | grep -v -i -E "none|white|rgb\s*\(\s*255\s*,\s*255\s*,\s*255\s*\)|rgb\s*\(255\s*255\s*255\s*\)|#f{3,6}|#e{2}f{2}f{2}|#d{2}f{2}f{2}|#b{2}f{2}f{2}|#a{2}f{2}f{2}")
        
        # Skip grayish colors (like #f0f0f0, #ececec, etc.)
        colors=$(echo "$colors" | grep -v -i -E "f2f2f2|f0f0f0|ececec|dcdcdc|bdbdbd|a9a9a9")

        # Only proceed if colors are found
        if [[ -n "$colors" ]]; then
            # Count the occurrences of each color
            most_common_color=$(echo "$colors" | sort | uniq -c | sort -nr | head -n 1 | awk '{print $2}')
            count=$(echo "$colors" | grep -c "$most_common_color")
            
            # Append the result for this file to the output JSON using jq
            output=$(echo "$output" | jq --arg svgName "$filename" --arg mostFrequentColor "$most_common_color" --argjson count "$count" \
                '. += [{"svgName": $svgName, "mostFrequentColor": $mostFrequentColor, "count": $count}]')
        fi
    fi

    # Increment the counter and calculate the progress
    counter=$((counter + 1))
    progress=$((100 * counter / total_files))

    # Print the progress bar
    printf "\rProgress: [%-50s] %d%%" "$(printf "#%.0s" $(seq 1 $((progress / 2))))" "$progress"
done

# Check if there are any results before saving the JSON
if [[ "$output" != "[]" ]]; then
    # Save the JSON to a file
    echo "$output" > most_frequent_colors.json
else
    echo "No valid colors found in any SVG files. No JSON file created."
fi

# Print a newline after the progress bar
echo -e "\nProcessing complete."
