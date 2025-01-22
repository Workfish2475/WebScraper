import json
from operator import contains


class DataFetcher:
    def __init__(self):
        self.data = None

    def parse(self, json_file):
        try:
            with open(json_file, 'r') as file:
                self.data = json.load(file)
            return
        except FileNotFoundError:
            raise FileNotFoundError(f"File {json_file} not found")
        except json.JSONDecodeError:
            raise json.JSONDecodeError("Invalid JSON format")

    def filterByName(self, name):
        if self.data is None:
            return None

        returnable_data = []
        name = name.lower()
        for data in self.data:
            item_name = data.get("name", "")

            if name in item_name.lower():
                returnable_data.append(data)

        return returnable_data

    def filterBySet(self, set):
        if self.data is None:
            return None

        returnable_data = []
        set = set.lower()
        for data in self.data:
            item_set = data.get("set", "")

            if set in item_set.lower():
                returnable_data.append(data)

        return returnable_data

    def filterByInfo(self, info):
        if self.data is None:
            return None

        returnable_data = []
        info = info.lower()
        for data in self.data:
            item_info = data.get("info", "")

            if set in item_info.lower():
                returnable_data.append(data)

        return returnable_data

    def filterByColor(self, color):
        if self.data is None:
            return None

        returnable_data = []
        color = color.lower()
        for data in self.data:
            item_color = data.get("color", "")

            if color in item_color.lower():
                returnable_data.append(data)

        return returnable_data
