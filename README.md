# Scraper.py

### How to run

- Run the command ```pip install -r requirements.txt```
- Once pip is done installing the needed dependencies, run ```python Scraper.py```

### What the script does

- The script will connect to a site and begin gathering all character/card information from the site. The script will navigate through each of pages available until there is none left.
- The script will then convert all the card information into a json file named ```data.json``` as shown below:

```
{
        "name": "Portgas.D.Ace",
        "cost": "7",
        "power": "7000",
        "counter": "-",
        "color": "Red",
        "type": "Whitebeard Pirates",
        "effect": "Effect[On Play] Give up to 2 of your opponent's Characters −3000 power during this turn. Then, if your Leader's type includes \"Whitebeard Pirates\", this Character gains [Rush] during this turn.(This card can attack on the turn in which it is played.)",
        "set": "Card Set(s)-TWO LEGENDS- [OP-08]",
        "attribute": "Special",
        "cardNo": 1,
        "imgPath": "assets/1.jpg",
        "info": "CHARACTER"
    },
```
- Once that step is complete, the script will then download each of the corresponding images from the website and store them in the ```assets/``` directory. 

### API Options
There are two options for launching the demo API: one in Python and the other in Go.

#### Python

Run the command: ```python server.py```

#### Go

Run the command: ```go run server.go```

#### Notes

- Currently, the APIs are just parsing the ```data.json``` file and returning everything based on that, but I plan to eventually make use of a DBMS to help make things a bit more efficient.