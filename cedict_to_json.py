# Read the file and process the data into the specified JSON format
file_path = "cedict_ts.u8"

# Dictionary to store the processed entries
entries = {}

# Processing the file
with open(file_path, "r", encoding="utf-8") as file:
    for line in file:
        # Skip comments and empty lines
        if line.startswith("#") or not line.strip():
            continue

        # Split the line into components
        parts = line.strip().split(" ")
        traditional, simplified = parts[0], parts[1]

        # Extracting pronunciation and definitions
        remaining = " ".join(parts[2:])
        pronunciation_start = remaining.find("[")
        pronunciation_end = remaining.find("]")
        pronunciation = remaining[pronunciation_start + 1 : pronunciation_end]

        definitions = remaining[pronunciation_end + 2 :].strip("/").split("/")

        # Process simplified entry
        if simplified not in entries:
            entries[simplified] = {
                "simplified": simplified,
                "traditional": traditional,
                "pronunciation": [pronunciation],
                "definitions": [definitions],
            }
        else:
            entries[simplified]["pronunciation"].append(pronunciation)
            entries[simplified]["definitions"].append(definitions)

        # Process traditional entry if different
        if traditional != simplified:
            if traditional not in entries:
                entries[traditional] = {
                    "simplified": simplified,
                    "traditional": traditional,
                    "pronunciation": [pronunciation],
                    "definitions": [definitions],
                }
            else:
                entries[traditional]["pronunciation"].append(pronunciation)
                entries[traditional]["definitions"].append(definitions)

# Convert the processed data to JSON
import json

json_data = json.dumps(entries, indent=4, ensure_ascii=False)

# Saving the output to a JSON file
output_file_path = "cedict.json"
with open(output_file_path, "w", encoding="utf-8") as output_file:
    output_file.write(json_data)
