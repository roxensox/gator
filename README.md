# Gator RSS Aggregator

## Description

Gator is a CLI RSS aggregator with a PostgreSQL backend. Intended usage is to run the ```agg``` command in the background persistently and interact with it in another terminal session. 

## Dependencies
- PostgreSQL
- Go
- MacOS or Linux

## Installation
In the root directory of this project, run the command 

```
go install .
```

Once this is complete, you should be able to run gator in your terminal emulator of choice.

### Configuration
Gator is configured by the .gatorconfig.json file in your home directory. The format is as follows:
```
{
    "db_url":"{PostgreSQL access string}",
    "current_user":"{LEAVE BLANK}"
}
```
You will also need to set up a psql server and connect the gator database to use this properly.

## Usage

### register
```gator register {username}``` allows you to register yourself as a user in Gator's database. As a registered user, you can curate which feeds you see when you run the ```browse``` command.

### login
```gator login {username}``` allows you to sign in as different users by username. This enables you to create different profiles for different interest if you follow a lot of feeds.

### addfeed
```gator addfeed {feed name} {feed url}``` allows you to add a feed to the database so that Gator can start tracking it.

### follow
```gator follow {feed url}``` allows you to follow a feed by URL when logged in on a profile.

### unfollow
```gator unfollow {feed url}``` allows you to remove followed feeds from your profile.

### agg
```gator agg {time interval}``` runs constantly in the background and collects posts from your rss feed at a chosen rate. Time formats are 1s (seconds), 1m (minutes), 1h (hours), 1d (days)

### reset
```gator reset``` allows you to fully reset the database.
