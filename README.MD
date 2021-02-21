### Matchmaker

When a player loads a game its game client sends a join request to the matchmaking service to let it know that the player wants to play the game and it awaits the game to start. As soon as there are 4 awaiting players the matchmaking service creates a new game session for them.
When a player closes a game its game client sends a leave request to the matchmaking service. If the player was in a game then the game ends for all other 3 players and they automatically start waiting for the next game to start.
There is a limit of 100 connected players per each deployed instance of the matchmaking service and we would like to have a /sessions endpoint which returns json list of running games with players ids per each game.

## Design

Matchmaking configs like the number of players allowed for a game session, maximum players for a deployed instance, etc... are loaded as environment variables during the server start.  
The environment variables are: MAX_PLAYERS_PER_INSTANCE and ALLOWED_PLAYERS_COUNT  
If no environment variables are provided then the default values are taken.

The matchmaker uses a channel to temporarily store the players who would like to join a game.  
On server start, the matchmaker is also started through a separate go routine.  
The matchmaker keeps reading from the channel and begins a new game session when the adequate number of players are waiting.

The game sessions and players data are stored in-memory using go's `sync.Map`, sync.Map is chosen to avoid race during parallel requests.


There are three endpoints to the server:
- /games/join
- /games/leave
- /games/session

# /games/join 

A player can use the `games/join` endpoint to join a new game session. A cookie is set in the response header. Subsequent requests to the same endpoint give the status of the player which could be either waiting or has successfully joined a game. If maximum players are already playing for the given instance then this requests would throw an error.


+ Response 202 (application/json)
```json
	{
	    "code": 202,
	    "message": "Successfully registered with id: e3af96c2-3c22-4c59-b94b-11fce20df7d4. Inadequate number of players to start the game. Please wait!"
	}
```

+ Response 200 (application/json)
```json
	{
	    "code": 200,
	    "message": "Successfully registered with id: e3af96c2-3c22-4c59-b94b-11fce20df7d4 and joined game with id: 4bb3f79d-fcdd-4ebc-9c29-a64200b8e9b6"
	}
```

+ Response 500 (application/json)
```json
	{
	    "code": 500,
	    "message": "Please wait, maximum players are playing in this instance."
	}
```


# /games/leave

A user can use the `games/leave` endpoint to leave a game session. Using the cookie that is set, the matchmaker would fetch the player details and end the corresponding game session. Other players in the game session would then be put into the wait channel.

+ Response 200 (application/json)
```json
	{
	    "code": 200,
	    "message": "Player has successfully left"
	}
```

+ Response 400 (application/json)
```json
	{
	    "code": 400,
	    "message": "Player is not playing any game!"
	}
```


# /games/session

Lists the number of active game sessions

+ Response 200 (application/json)
```json
	[
	    {
	        "id": "06ea47f8-9f99-486d-8ef4-61692516d07c",
	        "players": [
	            "7bf863de-efc8-4974-a9d0-79755c5730f5",
	            "7801dd35-c189-4d42-b1fb-89429aac8644",
	            "0b269828-d1d3-479d-a157-b7710de4f2e2"
	        ]
	    }
	]
```

The application also has extensive unit test coverage and CI pipeline.
We could use gRPC or HTTP2 to establish a web socket and avoid retaining cookies.

## Running the application

```
go run main.go
```

or

```
 docker build -t matchmaker .  

 docker run --publish 3001:3001 -i -t "matchmaker"       

 ```

## Pre-commit (Optional)
To ensure a clean code.

   Install pre-commit [https://pre-commit.com/#install]
   Install .pre-commit-config.yaml as a pre-commit hook

	pre-commit install

    Go static analysis tools run automatically on pre-commit. Run checks manually if needed using

	pre-commit run --all-files
