import json

from string import ascii_lowercase
from random import randint, choice

string_template = 'INSERT INTO `json_db`.`json_table` (`json_col`)' \
                    + ' VALUES (\'{json_val}\');'
json_template = {"name": "", "value": 0}

strings_array = []
for _ in range(100):
    name = ''.join(choice(ascii_lowercase) for _ in range(randint(5, 10))) \
        .capitalize()
    value = randint(10, 1000)
    json_template["name"], json_template["value"] = name, value
    new_json = json.dumps(json_template, separators=(",", ":"))
    insert_string = string_template.format(json_val=new_json)
    print(insert_string)

