import requests
import random
import time

random.seed(time.time())

def gen_char():
    # generate character randomly
    URL = "http://localhost:8080/newChar"

    classes = ['Knight','Archer','Wizard']
    c_player = classes[random.randint(0,2)]

    PARAMS = {
            'name': "Training",
            'class': c_player
            }

    requests.post(url = URL, params = PARAMS)

    c_enemy = classes[random.randint(0,2)]

    PARAMS = {
            'name': "Training_enemy",
            'class': c_enemy
            }

    requests.post(url = URL, params = PARAMS)

    return c_player, c_enemy

def play_game():
    c_player, c_enemy = gen_char()

    while True:
        moves = ['Heavy', 'Quick', 'Standard', 'Block', 'Parry', 'Evade']

        URL = ("http://localhost:8080/turn/Training." + c_player +
            "/Training_enemy." + c_enemy + "/" + moves[random.randint(0,5)])

        r = requests.get(url = URL)

        END_URL = "http://localhost:8080/end/"

        if END_URL in r.url:
            break

games_played = 0
while True:
    play_game()
    games_played += 1
    print(games_played)
