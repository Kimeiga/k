import json

ids_path = "ids.json"

# Load the ids JSON file
with open(ids_path, "r", encoding="utf-8") as file:
    ids_data = json.load(file)

new = {}

for char, ids_entry in ids_data.items():
    # if len(ids_entry) > 1:
    # print(len(ids_entry))
    new[char] = ids_entry[0]

new_json = json.dumps(new, indent=4, ensure_ascii=False)
with open("new.json", "w", encoding="utf-8") as output_file:
    output_file.write(new_json)
