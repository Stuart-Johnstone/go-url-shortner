# How am i going to build this url shortener?

## Challenges

There are going to be a couple of challenges. 

- CR database, append only for now
- Hash collisions, im not going to use sha265 for perf reasons
- cacheing
- make this work over http using a redirect, will need some home page probably


 DB Table format 
| Key | value |
| x6jg3p | google.com |
| 8blvjz | instagram.com/stuartkj |


## Steps

### Step 1: Learn the basics of go
- take in usr input, create a hash of the user input
-- Think of this as creating the shortened url

### Step 2: Learn how to setup the mongodb to interact with user input
- learn to write to the db (hash as key, url as value)
- learn how to read from the db

### Step 3: Integrate the go server with http calls to localhost
- some kind of console.log equivalent 
- redirect urls

### Step 4: Add a redis cache for better performance
- this one is performative for the size of the db but would be cool to have

### Step 5: Containerize for fun
