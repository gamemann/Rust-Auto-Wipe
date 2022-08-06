# Rust Auto Wipe (WIP/In Development)
## Description
An application made in Go for Rust servers operating with [Pterodactyl](https://pterodactyl.io/). This application automatically wipes server(s) weekly by default, but supports options such as monthly and biweekly. The program is aimed to be as flexible as possible. With that said, there are many features including the following.

* Allow rotating of map seeds.
* Allow automatically changing the host name on each wipe with format support (including option replacements like `{day}` and `{month}`).
* Deletion of files with the option to except specific types (e.g. don't delete player data such as states, identities, tokens, and more).
* A flexible configuration and timing system.
* Support for retrieving servers from Pterodactyl API and allowing environmental overrides.

## Configuration
All configuration is done inside a config file on disk via JSON. The default path is `/etc/raw/raw.conf` and may be changed via the `-cfg` flag.

Any wipe-specific configuration at the top-level of the configuration is used as the default values for each server. Each server may override these by specifiying the same key => value pairs inside of the server array. 

There is also support for environmental variables for auto-added servers from the Pterodactyl panel which will allow you to give access to override server-specific options to users with authorization in your Pterodactyl database.

When the program is ran, but no configuration file is found, it will attempt to create the file (likely requiring root privileges by default if trying to create inside of `/etc/`, however). The below are all JSON settings along with their default values, but with added comments. Remember that JSON does **not** support comments. Therefore, copying the below contents with comments will result in errors. This is why it's recommended to allow the program to create a default configuration file.

```json
{
    // The URL to the panel (make sure to include a trailing /). Ex: http://ptero.something.internal/
    "apiurl": "",

    // The Pterodactyl client token (should start with "ptlc_"). Create under user account settings in Pterodactyl panel.
    "apitoken": "",

    // Debug level from 0 - 4.
    "debuglevel": 1,

    // The application token (required for automatically adding servers from Pterodactyl).
    "apptoken": "",

    // Automatically add servers from Pterodactyl (servers require 'WORLD_SEED' and 'HOSTNAME' environmental variables).
    "autoaddservers": false,

    // Path starting from /home/container to the server files we need to delete (should be /server/rust with default Rust egg).
    "pathtoserverfiles": "/server/rust",

    // Timezone from Go's time library.
    "timezone": "America/Chicago",

    // The weekly wipe time (<day> <hour>:<minute>). Hour is from 0 - 23 (using 24-hour system) and minute is 0 - 60.
    "wipetime": "Thursday 12:00",

    // Special flag to wipe on the first of each month.
    "wipemonthly": false,

    // Special flag to wipe biweekly.
    "wipebiweekly": false,

    // Whether to delete map files (includes *.map and *.sav files).
    "deletemap": true,

    // Whether to delete player blueprints (includes any files with blueprints in the file name).
    "deletebp": true,

    // Whether to delete deaths (includes any files with deaths in the file name).
    "deletedeaths": true,

    // Whether to delete states (includes any files with states in the file name). 
    "deletestates": true,

    // Whether to delete identities (includes any files with identities in the file name). 
    "deleteidentities": true,

    // Whether to delete tokens (includes any files with tokens in the file name). 
    "deletetokens": true,

    // Whether to delete server data/files (includes any files with sv.files in the file name). 
    "deletesv": true,

    // Whether to change the map seed.
    "changemapseed": false,

    // A list of map seeds (e.g. [1203213, 12312312, 235123]).
    "mapseeds": null,

    // Pick type (1 = pick the next seed, otherwise, pick a random seed).
    "mapspicktype": 1,

    // Whether to change the hostname.
    "changehostname": true,

    // The hostname format.
    // Replacements include:
    // {month} = 1 - 12.
    // {day} = 1 - x.
    // {week_day} = 0 - 6 (starting from Sunday).
    // {week_day_str} = Monday - Sunday
    // {seconds_left} = Used for warning messages, but how many seconds left until wipe.
    "hostname": "Vanilla | FULL WIPE {month}/{day}",
    
    // Whether to merge both server-specific and global warning messages.
    "mergewarnings": false,

    // Warning messages list
    "warningmessages": [
        {
            // The prewarn time (e.g. this would warn one second before wipe whereas if the warning time was 10, it would warn 10 seconds before wipe time).
            "warningtime": 1,

            // The message (use formatting from hostname documentation).
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 2,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 3,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 4,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 5,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 6,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 7,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 8,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 9,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        },
        {
            "warningtime": 10,
            "message": "Wiping server in {seconds_left} seconds. Please join back!"
        }
    ],

    // Server list (null by default).
   "servers": null
}
```

The servers array includes the following:

```json
{
    "servers": [
        {
            // Whether to enable the server or not (enabled by default).
            "enabled": true,

            // The (short) UUID of the server. Characters before the first "-" in the long UUID.
            "uuid": "",

            // Overrides (retrieve definition from top-level comments above).
            "apiurl": "",
            "apitoken": "",
            "debuglevel": 1,
            "pathtoserverfiles": "/server/rust",
            "timezone": "America/Chicago",
            "wipetime": "Thursday 12:00",
            "wipemonthly": false,
            "wipebiweekly": false,
            "deletemap": true,
            "deletebp": true,
            "deletedeaths": true,
            "deletestates": true,
            "deleteidentities": true,
            "deletetokens": true,
            "deletesv": true,
            "changemapseeds": false,
            "mapseeds": null,
            "mapspicktype": 1,
            "changehostname": true,
            "hostname": "Vanilla | FULL WIPE {month}/{day}",
            "mergewarnings": false,
            "warningmessages": null,

            // Extras.
            // Wipe server when the program is first started.
            "wipefirst": false
        }
        ...
    ]
}
```

**Note** - When writing to the default file after creation, it will try to make the JSON data pretty (AKA pretty print by idents).

## Environmental Overrides With Auto-Added Servers
There are environmental overrides for servers that are added from the Pterodactyl API. This allows you to distribute access easier from within the Pterodactyl panel itself.

These are only server-specific options for obvious reasons.

It doesn't technically matter what type of variable you make these (they are parsed properly within the code). However, it's recommended to still use integers, booleans, and so on for readability reasons instead of strings.

The following is a list of environmental names you can create variables within Pterodactyl Nests/Eggs for overrides.
* **RAW_ENABLED** - Enabled override.
* **RAW_PATHTOFILES** - Path to files override.
* **RAW_TIMEZONE** - Timezone override.
* **RAW_WIPETIME** - Wipe time override.
* **RAW_WIPEMONTHLY** - Wipe monthly override.
* **RAW_WIPEBIWEEKLY** - Wipe biweekly override.
* **RAW_DELETEMAP** - Delete map override.
* **RAW_DELETEBP** - Delete blueprints override.
* **RAW_DELETEDEATHS** - Delete deaths override.
* **RAW_DELETESTATES** - Delete states override.
* **RAW_DELETEIDENTITIES** - Delete identities override.
* **RAW_DELETESV** - Delete server files/data override.
* **RAW_CHANGEMAPSEEDS** - Change map seeds override.
* **RAW_MAPSEEDS** - Map seeds override (this is a special case, map seeds should be a string with integers representing the seed list separated by commas, ",". Example - "123213,12314,123412").
* **RAW_MAPSEEDSPICKTYPE** - Change map seeds pick type override.
* **RAW_MAPSEEDSMERGE** - Change map seeds merge override.
* **RAW_CHANGEHOSTNAME** - Change hostname override.
* **RAW_HOSTNAME** - Hostname override.
* **RAW_MERGEWARNINGS** - Merge warnings override.
* **RAW_WARNINGMESSAGES** - Warning messages override (another special case, this should be a JSON string of the normal `warningmessages` JSON item). Example - `{"warningmessages": [{"warningtime": 5, "message": "{seconds_left} until wipe!"}]}`.
* **RAW_WIPEFIRST** - Wipe first override.

## Credits
* [Christian Deacon](https://github.com/gamemann)