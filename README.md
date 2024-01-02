# React Roles

React Roles is a simple discord bot for allowing self-managing user roles in discord.

<p align="center">
    <img src="https://user-images.githubusercontent.com/26305909/170378884-1969ed52-799a-4387-9beb-8187859c9750.png" alt="Role message with roles"/>
    <img src="https://github.com/Zaptross/reactroles/assets/26305909/9643bd06-6d91-40e2-b297-2beac6cffdca" alt="Slash commands interface"/>
</p>

## Usage

To add a role to yourself as a user you simply click one of the reactions on the message, and to remove yourself from that role you click again to remove your reaction on the message.

If a user has the roles (as configured below) that gives the permission to use role management commands:

### Conventions

- <> : a required parameter
- [] : an optional parameter

### Configuration

The bot is configured per server via the `configure` command.

To configure the bot, you will need:

- A channel for the bot to post the role selector messages in (e.g. `#roles`)
- One or more roles to use as permission roles for managing roles (e.g. `@role-add, @role-update, @role-remove`)
  - These can be the same role if you want to give all permissions to one role
- Is creating and removing channels enabled? (e.g. `true, false`)
- What category should the bot create new role channels in? (e.g. `📁 role-channels`)
  - This is only used if creating and removing channels is enabled
- One or more roles to use as permission roles for managing role channels (e.g. `@role-channel-create, @role-channel-remove`)
  - These can be the same as the role management role

#### Example

- `/role configure <role channel> <add role> <remove role> <update role> <create channel> <create channel role> <remove channel role> <channel category> <cascadeDelete>`
- `/role configure #roles @role-add @role-remove @role-update true  @role-channel-create @role-channel-remove role-channels true`

### Add Command

Adds a new role to the discord, configured as specified.

- `/role add <role name> <emoji> [colour]`

Usage examples:

- `/role add valorant :gun: #d34454`
- `/role add valorant :gun:`

### Remove Command

Removes a role, its reactions, and its channels if they exist.

- `/role remove <role>`

Usage example:

- `/role remove @valorant`

### Update Command

Modifies any one part of a role.
Where role fields are `name`, `emoji` and `color`, and role field values are valid values of those fields as per the `add` role command.

- `/role update <role> <role field> <role field value>`

Usage examples:

- `/role update @valorant name coolgungame`
- `/role update @valorant emoji 😎`
- `/role update @valorant color #CADEAA`

### Create Channel Command

Creates a new channel with the specified name, and adds the specified role to the channel.
Roles may have zero or one text channel, and zero or one voice channel associated with them.

- `/role create-channel <role> <channel name> <channel type>`

Usage examples:

- `/role create-channel @valorant valorant-chat text`
- `/role create-channel @valorant valorant-voice voice`

### Link Channel Command

Links an existing channel to a role.
Roles may have zero or one text channel, and zero or one voice channel associated with them.

- `/role link-channel <role> <channel>`

Usage examples:

- `/role link-channel @valorant #valorant-chat`
- `/role link-channel @valorant 🔊valorant-voice`

### Remove Channel Command

Removes a channel from a role, and deletes the channel.

- `/role remove-channel <role> <channel>`

Usage examples:

- `/role remove-channel @valorant #valorant-chat`

### Help Command

Replies to the user the help text accompanying the command.

- `/role help <action>`

Usage example:

- `/role help add`

# Setup

Before deploying your own, you will need to make a discord bot, and add it to your server.

- Discord application creation: https://discord.com/developers/applications
- Discord oauth2 link generator(with correct permissions preconfigured): https://discordapi.com/permissions.html#2415995968

## Docker

You can check out the image versions over on [Docker Hub](https://hub.docker.com/r/zaptross/reactroles)

1. Clone the repo to your machine
2. Duplicate `./example.env` and rename it to `.env`
3. Fill out the env variables
4. In a terminal in the repo root, run `docker-compose up -d`

## Kubernetes

1. Clone the repo to your machine
2. Duplicate `./cmd/reactroles/kube/deployment.yaml` and rename it to `deployment.prod.yaml`
3. Fill out the data variables in the config map
4. In a terminal in the repo root, run `kubectl apply -f ./cmd/reactroles/kube/deployment.prod.yaml`
