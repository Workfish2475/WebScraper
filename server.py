from flask import Flask
from jsonParse import DataFetcher
import json

app = Flask(__name__)

@app.route('/fetchAll')
def fetchall():
    fetcher.parse('data.json')
    results = fetcher.filterByName("")
    return results

@app.route('/fetchMatchesByName/<string:name>')
def fetchMatchesByName(name: str):
    fetcher.parse('data.json')
    results = fetcher.filterByName(name)
    return results

@app.route('/fetchMatchesBySet/<string:set>')
def fetchMatchesBySet(set: str):
    fetcher.parse('data.json')
    results = fetcher.filterBySet(set)
    return results

@app.route('/fetchMatchesByColor/<string:color>')
def fetchMatchesByColor(color: str):
    fetcher.parse('data.json')
    results = fetcher.filterByColor(color)
    return results

@app.route('/fetchMatchesByInfo/<string:info>')
def fetchMatchesByInfo(info: str):
    fetcher.parse('data.json')
    results = fetcher.filterByInfo(info)
    return results

if __name__ == "__main__":
    fetcher = DataFetcher()
    app.run(debug=True)
