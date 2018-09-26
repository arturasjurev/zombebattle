[![Build Status](https://travis-ci.com/sheirys/zombebattle.svg?branch=master)](https://travis-ci.com/sheirys/zombebattle)
[![Go Report Card](https://goreportcard.com/badge/github.com/sheirys/zombebattle)](https://goreportcard.com/report/github.com/sheirys/zombebattle)
[![codecov](https://codecov.io/gh/sheirys/zombebattle/branch/master/graph/badge.svg)](https://codecov.io/gh/sheirys/zombebattle)
[![GoDoc](https://godoc.org/github.com/sheirys/zombebattle/engine?status.svg)](https://godoc.org/github.com/sheirys/zombebattle/engine)

`go get -d github.com/sheirys/zombebattle/engine`

# zombebattle

zombebattle is simple framework to implement game based on mysterium network [Communication channel specification](https://github.com/mysteriumnetwork/winter-is-coming/blob/master/quests/Talk_to_Zombies.md#communication-channel-specification). **This project should be considered as framework to implement game based on this specification**. You will find some example server implementations in `engine/examples`.

## TrainingGrounds 
In `engine/examples/demo/main.go` you can find `TrainingGrounds` game implementation. This implementation will use one room where no-one can win. You will find stationary `dummy` zombie there, that is always staying in same position `x:5, y:5` and has unlimited HP (cannot be killed). This implementation supports multiple users in same room. To start use `go run engine/examples/demo/main.go`. This will start tcp server on `:3333`. You can connect clients with `telnet localhost 3333`. Type `start <name>` to join the game (e.g. `start warrion1`).

You can try to shoot dummy as client with `shoot 5 5` command.

## TheWall
In `engine/examples/thewall/main.go` you will find `TheWall` game implementation. This implementation will spawn zombie when new client joins the room. Zombies will try to reach the wall, and if they reach wall 5 times, zombies will win. You must shoot 5 zombies to win this room. Score count cant be set at `engine/rooms/TheWall.go:TheWallMaxPlayerScore` or `engine/rooms/TheWall.go:TheWallMaxZombieScore`.

## MutliClient support
This implementation support multiple clients. Clients can join same room and kill zombies together. Communication channel specification requires that if multiple clients are joined same room, first good shot will end the game. IMHO this is not fun, so if you want to win this round players need to score 5 points (kill 5 zombies). However this requirement can be implemented by changing `engine/rooms/TheWall.go:TheWallMaxPlayerScore` to `1`, so players will win when they score 1 point.

## MultiRoom support
As Communication channel specification by default does not specify how to support multiple rooms with multiple clients at one time (no instuctions how to two players can play two separate games), this implementation has extended Communication channel specification with additional commands `JOIN` and `NEW`. When client is connected into lobby (connected by telnet but not executed `START` command), client can create new rooms with `NEW <name>` command (e.g. `new world1`) or select room where he wants to join `JOIN <name>`. If client does not select the room with `JOIN` command after `START` he will be forced to join default room.

Client usage example for single room:

        # telnet localhost 3333
        start vanagas    # will join default room as vanagas
        shoot 1 1        # try to shoot zombie

Client usage example for multi room:

        # telnet localhost 3333
        new world1       # create new room
        join world1      # select to enter this room
        start vanagas    # will join world1 as vanagas
        shoot 1 1        # try to shoot zombie

Example server usage:
```
	server := &engine.Server{
		Addr: ":3333",
                // predefine default room. here clients will be forced to join
                // if didn't selected room with `JOIN` command.
		DefaultRoom: &rooms.TrainingGrounds{
                        // this room has predefined dummy zombie
			Zombies: []types.Zombie{
				&zombies.Dummy{},
			},
		},
                // predefine other rooms. server is not limited to these rooms.
                // client can create new room with `NEW` command.
		Rooms: []types.ServerRoom{
			{Room: &rooms.TheWall{}, Default: false},
		},
	}

	server.Run()
```
