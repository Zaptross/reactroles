# React Roles

React Roles is a simple discord bot for allowing self-managing user roles in discord.

<p align="center">
    <img src="https://user-images.githubusercontent.com/26305909/170378884-1969ed52-799a-4387-9beb-8187859c9750.png" alt="Role message with roles"/>
</p>

## Usage

To add a role to yourself as a user you simply click one of the reactions on the message, and to remove yourself from that role you click again to remove your reaction on the message.

If a user has the roles (as configured below) that gives the permission to use role management commands:

-   `!role <add/remove> <role name> <emoji> [colour]`

Usage examples:

-   `!role add valorant :gun: #d34454`
-   `!role add valorant :gun:`
-   `!role remove valorant`

When a command is successful, the bot will remove the your message containing the command and update the roles selector accordingly.

If unsuccessful, the bot will reply with an error message and usage then delete your message containing the command.

# Setup

Before deploying your own, you will need to make a discord bot, and add it to your server.

-   Discord application creation: https://discord.com/developers/applications
-   Discord oauth2 link generator(with correct permissions preconfigured): https://discordapi.com/permissions.html#268512320

## Docker

1. Clone the repo to your machine
2. Duplicate `./example.env` and rename it to `.env`
3. Fill out the env variables
4. In a terminal in the repo root, run `docker-compose up -d`

## Kubernetes

1. Clone the repo to your machine
2. Duplicate `./cmd/reactroles/kube/deployment.yaml` and rename it to `deployment.prod.yaml`
3. Fill out the data variables in the config map
4. In a terminal in the repo root, run `kubectl apply -f ./cmd/reactroles/kube/deployment.prod.yaml`
