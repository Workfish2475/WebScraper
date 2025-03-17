from flask import Flask, jsonify, render_template
from jsonParse import DataFetcher
import json

app = Flask(__name__)
fetcher = DataFetcher()


@app.before_request
def initialize():
    fetcher.parse('data.json')


@app.route('/')
def home():
    return render_template('index.html')

@app.route('/fetchAll')
def fetchall():
    results = fetcher.getData()
    return jsonify(results)


@app.route('/name/<string:name>')
def fetchMatchesByName(name: str):
    results = fetcher.filterByName(name)
    return jsonify(results)


@app.route('/set/<string:set>')
def fetchMatchesBySet(set: str):
    results = fetcher.filterBySet(set)
    return jsonify(results)


@app.route('/color/<string:color>')
def fetchMatchesByColor(color: str):
    results = fetcher.filterByColor(color)
    return jsonify(results)


@app.route('/info/<string:info>')
def fetchMatchesByInfo(info: str):
    results = fetcher.filterByInfo(info)
    return jsonify(results)


if __name__ == "__main__":
    fetcher.parse('data.json')
    app.run(debug=True)
