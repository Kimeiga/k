# Convert the processed data to JSON
import json

# Path to the cedict JSON file

cedict_path = "cedict.json"

# Load the cedict JSON file
with open(cedict_path, "r", encoding="utf-8") as file:
    cedict_data = json.load(file)

# Path to the ids JSON file
ids_path = "ids.json"

# Load the ids JSON file
with open(ids_path, "r", encoding="utf-8") as file:
    ids_data = json.load(file)

# Combine the data
combined_data = {}

for char, cedict_entry in cedict_data.items():
    combined_data[char] = {
        "cedict": cedict_entry,
        "ids": ids_data.get(char, {}),  # Get the ids entry if it exists
    }

# Also include characters from ids that are not in cedict
for char, ids_entry in ids_data.items():
    if char not in combined_data:
        combined_data[char] = {"cedict": {}, "ids": ids_entry}

# Convert the combined data to JSON
combined_json = json.dumps(combined_data, indent=4, ensure_ascii=False)

# Save the combined data to a JSON file
output_file_path = "combined.json"
with open(output_file_path, "w", encoding="utf-8") as output_file:
    output_file.write(combined_json)
