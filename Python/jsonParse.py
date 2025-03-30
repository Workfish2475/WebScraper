import json

class DataFetcher:
    def __init__(self):
        self.data = None

    def parse(self, json_file):
        try:
            with open(json_file, 'r') as file:
                self.data = json.load(file)
        except FileNotFoundError:
            raise FileNotFoundError(f"File {json_file} not found")
        except json.JSONDecodeError:
            raise json.JSONDecodeError("Invalid JSON format")

    def _filter_by_field(self, field, value):
        if self.data is None:
            return None

        value = value.lower()
        return [item for item in self.data if value in (item.get(field, "")).lower()]

    def filterByName(self, name):
        return self._filter_by_field("name", name)

    def filterBySet(self, set_name):
        return self._filter_by_field("set", set_name)

    def filterByInfo(self, info):
        return self._filter_by_field("info", info)

    def filterByColor(self, color):
        return self._filter_by_field("color", color)

    def getData(self):
        return self.data