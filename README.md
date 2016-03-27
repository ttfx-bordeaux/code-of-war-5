# Code of War : Enlarge your tower ![](https://travis-ci.org/ttfx-bordeaux/code-of-war-5.svg?branch=master)
![](http://www.codeofwar.net/sites/all/themes/cow/images/xlogo.png.pagespeed.ic.n8tK1fUftd.png)

## Staff :
- [Target process](https://kriyss.tpondemand.com)
- [Travis](https://travis-ci.org/ttfx-bordeaux/code-of-war-5)
- [google-doc](https://docs.google.com/document/d/1mAcHqqwybe-Z9JYzGX4Fi2q3ZZmjIFUjllQGPF7tQ-w/edit?usp=sharing)

## Server

### How to launch ?

Default `game` server port value is `3000`.  
Default `admin` server port value is `4000`
```sh
  go build && ./server [--port <port>] [--admin-port <admin-port>]
```



## Message on Game Server 

### Authentication

```json 
{
    "action" : "authenticate",
    "data": {
        "id":"token",
        "name":"username"
    }
 }
```

## Message on Admin Server

### Start Game

route                                   |method  |   action
----------------------------------------|--------|-----------------------------------
/game                                   | POST   |  create new game, return `game_id`
/game/`{game_id}`/player/`{player_id}`  | POST   |  add player to game
/game/`{game_id}`/player/`{player_id}`  | DELETE |  remove player from game
/game/`{game_id}`/launch                | POST   |  launch the game

