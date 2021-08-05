# Skelly

Skelly is a Slack bot for automatically reacting to user typing, written in [Go](https://golang.org/)

## Usage

To configure Skelly, all you need to do is

1. Invite Skelly to your Slack channel

```
/invite @skelly
```

2. Add a reaction in [Slack](https://slack.com/) through Skelly's [slash commands](https://slack.com/help/articles/201259356-Use-built-in-slash-commands)

[See also](https://api.slack.com/interactivity/slash-commands)

```
/skelly help
/skelly add @slack-admins
```

3. Start typing!

If you want your response to include nice links, use the following syntax:

```none
<http://www.foo.com|This message is a link>
```

## Managing Reactions

For a list of useful commands, try

```
/skelly help
```

| Command  | Input | Effect |
| ------------------- | -------- | ------------- |
| /skelly help  | NONE | prints helpful information |
| /skelly add  | NONE | opens the modal for adding a typing reaction _in that channel_ |
| /skelly add | NONE | adds a typing reaction _in that channel_ for all users |
| /skelly update | NONE | opens the modal for updating a typing reaction |
| /skelly delete | NONE | opens the modal for deleting a typing reaction |
| /skelly list  | NONE` | lists all typing reactions that exist _in that channel_ |



## Development

To run the bot locally, simply configure the environment and use the `Makefile`

Interact with Skelly using the built-in CLI to make developing locally easier.

```bash

$ make build

$ ./release/skelly --help

$ ./release/skelly reaction add --channel <CHANNEL_ID> --response "Hello!"

$ ./release/skelly reaction trigger --channel <CHANNEL_ID> --user <USER_ID>

```

### Environment

Store the required configurations in either your environment or an `.env` file

| Variable  | Source |
| ------------- | ------------- |
| SKELLY_BOT_TOKEN  | [Slack bot token](https://api.slack.com/authentication/token-types#granular_bot) |
| SKELLY_VERIFICATION_TOKEN  | [Slack verification token](https://api.slack.com/authentication/verifying-requests-from-slack) |
| SKELLY_SIGNING_SECRET | [Slack signing secret](https://api.slack.com/authentication/verifying-requests-from-slack) |
| SKELLY_MONGO_HOST | [Mongo DB host](https://docs.mongodb.com/manual/reference/program/mongo/) |
| SKELLY_MONGO_DB | [Mongo DB database name](https://docs.mongodb.com/manual/reference/program/mongo/) |
| SKELLY_MONGO_USERNAME | [Mongo DB username](https://docs.mongodb.com/manual/tutorial/enable-authentication/) |
| SKELLY_MONGO_PASSWORD | [Mongo DB password](https://docs.mongodb.com/manual/tutorial/enable-authentication/) |

### Make

Use the `Makefile` to build and run the binary or the Docker image

```bash
# clone and navigate to skelly
$ git clone git@github.com:davidvader/skelly.git
$ cd skelly

# build the skelly binary and run the server
$ make up
```

Trigger the bot with 
You can also simulate all of Skelly's bot functionality using the Skelly CLI

```bash
# build the skelly binary
$ make build

# move the binary to bin
$ cp release/skelly /usr/local/bin/

# run the application
$ skelly --help

# add a reaction
$ skelly reaction view --channel C016DRZPLBC

# list reactions
$ skelly reaction list --channel C016DRZPLBC
```
