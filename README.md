# Go Aggregator CLI App

This is a project I created following the guided project on [Boot.dev](https://www.boot.dev). It's an RSS aggregator that works as a CLI app. There is no web server, just a command line tool but a postgres database is required. This is a go cli app so you will need go and it's toolchain installed to install this cli app.

## Setup Go

Make sure you have at least Go 1.25.6 installed from the [go install page](https://go.dev/doc/install) directly.

## Setup Postgresql

Postgresql is needed to keep track of the rss entries. You can find an installation method for your OS on the [postgres installation page](https://www.postgresql.org/download/).

## Configure Aggregator

The cli app requires a config file in your home directory called `.gatorconfig.json`. This config file should have the following fields:

```json
{
  "db_url": "postgres://username:@hostname:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

The postgres connection string can be any valid postgres connection string, see this [artcile about postgresql connection strings](https://www.geeksforgeeks.org/postgresql/postgresql-connection-string/) for if you need a guide to what makes a valid connection string. If you are using a local postgresql installation, that `sslmode=disable` parameter is recommended.

Make sure your database has a database called `gator` by entering `psql` and running:

```sql
CREATE DATABASE gator;
```

## Database migrations

Your database needs to be migrated to ensure your database works with this app. Make sure to git clone this project somewhere on your computer. Once that is done, you need to install [goose](https://github.com/pressly/goose), a Go based database migration tool.

Now that goose is installed, you can use the goose cli tool to run the migration. All you need is the same postgres connection string and to `cd` into the directory where your migrations live. It will look something like this:

### How to get your db connection string:

```bash
cat ~./gatorconfig.json
```

You can omit the parameters in your connection string when you use it with goose.

### Run the database migration

Assuming you're in the root of this projects repo on your system, execute the following:

```bash
cd sql/schema;

# now in sql/schema run this goose command
goose postgres '<postgres connection string>' up
```

If this command is successful, then you can proceed. If you encounter an error, you can visit the [goose repo](https://github.com/pressly/goose) to troubleshoot.

## Install aggregator

While you can run `go run .` from the root of this project in your local copy of the repo, you can install it by simply running `go install .`

Now you will be able to use the CLI app anywhere in your terminal

```bash
aggregator register '<your name>'
```

## How to use

Aggregator has the following commands:

- register
- login
- users
- addfeed
- feeds
- follow
- unfollow
- following
- agg
- browse
- reset

A general overview of how this CLI app works is like this. You register a user in your database: `aggregator register john`, this will login automatically (save the current user in the `.gatorconfig.json`). You can add RSS feeds: `aggregator addfeed 'Hacker News' '<url>'` this will save the feeds in the database.

You can see the feeds in your database by using the `aggregator feeds` command. If a RSS feed has already been added by another registered user, you can just follow that feed without using the `addfeed` command: `aggregator follow <url>`. You can find that url by using the `aggregator feeds` command. You can also unfollow a feed similarly with `aggregator unfollow <url>`.

To start collecting posts, run the `aggregator agg <time duration>` command. time duration is a value like `10s`, `30s`, `15m`, or `1h`. For example `aggregator agg 30s` will fetch one feed every 30 seconds, fetching the feed that hasn't been fetched the longest.

Finally, to see the posts, use `aggregator browse <limit>`. limit is the maximum lastest post you want to see in your terminal. The default is 2.

### Command: register

This command registers a user in your database:

```bash
# aggregator register <name>

aggregator register john
```

When you register a new user, the current user is saved to the `.gatorconfig.json` in your home directory.

### Command: login

This command can switch between users of the CLI app that have already been registered. This command updates the current user in the `.gatorconfig.json` in your home directory.

```bash
# aggregator login <name>

aggregator login kevin
```

### Command: users

This command displays a list of all registered users and shows who is the currently logged in user.

```bash
aggregator users
```

### Command: addfeed

This command adds a RSS feed to the database. a feed can only be added once by url, but can be followed by any registered user.

```bash
# aggregator addfeed <name> <url>

aggregator addfeed 'Hacker News' https://news.ycombinator.com/rss
```

### Command: feeds

This shows information on all the feeds saved to the database.

```bash
aggregator feeds
```

### Command: follow

This command allows the currently logged in user to follow a specific feed saved in the database even if it was added by another user.

```bash
# aggregator follow <url>

aggregator follow https://news.ycombinator.com/rss
```

Keep in mind that if you added this feed, it automatically gets followed for your user.

### Command: unfollow

As it sounds like, it allow the current user to unfollow a feed.

```bash
# aggregator unfollow <url>

aggregator unfollow https://news.ycombinator.com/rss
```

### Command: following

Lists the feeds the current user is following.

```bash
aggregator following
```

### Command: agg

This command starts the loop to fetch the oldest feeds that need to be refetched one by one according to the time between requests specified. new posts from these feeds that are not saved will be saved to the database for later browsing.

```bash
# aggregator agg <time between requests>

aggregator agg 30s
```

The time between requests follow this format: `10s`, `30s`, `15m`, `1h`.

### Command: browse

This command lets you browse the posts saved when the feeds where fetched. You can specify how many of the recent posts you want displayed to the terminal.

```bash
# aggregator browse <limit>

aggregator browse 5
```

### Command: reset

this command deletes all registered users, which cascades to deleting all feeds, and saved posts. Essentially giving you a blank slate. **CAREFUL:** This will mean you will have to re-add and follow all the feeds previously added.

```bash
aggregator reset
```
