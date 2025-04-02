import json
import mysql.connector

conn = mysql.connector.connect(
    host="placeholder",
    user="placeholder",
    password="placeholder",
    database="placeholder"
)

cursor = conn.cursor()

with open("../data.json", "r") as file:
    data = json.load(file)

for card in data:
    sql = """
    INSERT INTO cards (Name, Cost, Power, Counter, Color, Type, Effect, CardSet, Attribute, CardNo, ImgPath, Info)
    VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
    """
    values = (
        card["name"], card["cost"], card["power"], card["counter"], card["color"],
        card["type"], card["effect"], card["set"], card["attribute"], card["cardNo"], card["imgPath"], card["info"]
    )
    cursor.execute(sql, values)

conn.commit()
cursor.close()
conn.close()
