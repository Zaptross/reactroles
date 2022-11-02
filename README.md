# React Roles

React Roles is a simple discord bot for allowing self-managing user roles in discord.

<p align="center">
    <img src="https://user-images.githubusercontent.com/26305909/170378884-1969ed52-799a-4387-9beb-8187859c9750.png" alt="Role message with roles"/>
</p>

## Usage

To add a role to yourself as a user you simply click one of the reactions on the message, and to remove yourself from that role you click again to remove your reaction on the message.

If a user has the roles (as configured below) that gives the permission to use role management commands:

### Conventions

- < > : a required parameter
- [ ] : an optional parameter

### Add Command

Adds a new role to the discord, configured as specified.

- `!role add <role name> <emoji> [colour]`

Usage examples:

- `!role add valorant :gun: #d34454`
- `!role add valorant :gun:`

### Remove Command

Removes a role and it's reacions from the discord.

- `!role remove <role name>`

Usage example:

- `!role remove valorant`

### Update Command

Modifies any one part of a role.
Where role fields are `name`, `emoji` and `color`, and role field values are valid values of those fields as per the `add` role command.

- `!role update <role name> <role field> <role field value>`

Usage examples:

- `!role update valorant name coolgungame`
- `!role update valorant emoji 😎`
- `!role update valorant color #CADEAA`

### Help Command

Replies to the user the help text accompanying the command.

- `!role help <action>`

Usage example:

- `!role help add`

As of v2.6.0 (3ee3b59)

If unsuccessful, the bot will reply with an error message and usage.

# Setup

Before deploying your own, you will need to make a discord bot, and add it to your server.

- Discord application creation: https://discord.com/developers/applications
- Discord oauth2 link generator(with correct permissions preconfigured): https://discordapi.com/permissions.html#268512320

## Docker

You can check out the image versions over on [Docker Hub](https://hub.docker.com/repository/docker/zaptross/reactroles)

1. Clone the repo to your machine
2. Duplicate `./example.env` and rename it to `.env`
3. Fill out the env variables
4. In a terminal in the repo root, run `docker-compose up -d`

## Kubernetes

1. Clone the repo to your machine
2. Duplicate `./cmd/reactroles/kube/deployment.yaml` and rename it to `deployment.prod.yaml`
3. Fill out the data variables in the config map
4. In a terminal in the repo root, run `kubectl apply -f ./cmd/reactroles/kube/deployment.prod.yaml`
