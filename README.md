# Rust Auto Wipe
[![RAW Build Workflow](https://github.com/gamemann/Rust-Auto-Wipe/actions/workflows/build.yml/badge.svg)](https://github.com/gamemann/Rust-Auto-Wipe/actions/workflows/build.yml) [![RAW Run Workflow](https://github.com/gamemann/Rust-Auto-Wipe/actions/workflows/run.yml/badge.svg)](https://github.com/gamemann/Rust-Auto-Wipe/actions/workflows/run.yml)

## Description
An application made in Go (latest version works) for Rust servers operating with [Pterodactyl](https://pterodactyl.io/) (latest API version works). This application automatically wipes server(s) based off of cron jobs. The program is aimed to be as flexible as possible. With that said, there are many features including the following.

* Allow rotating of world settings (map/level, world size, and world seed).
* Allow automatically changing the host name on each wipe with format support (including option replacements like `{day_two}` and `{month_two}`).
* Deletion of files with the option to except specific types (e.g. don't delete player data such as states, identities, tokens, and more).
* A flexible configuration and uses a cron job system (support for multiple cron jobs per server).
* Support for retrieving servers automatically from Pterodactyl API.
* Pterodactyl egg variable overrides to allow others to change these settings without having access to the application's file system.
* Pre and post wipe hooks (sends POST data to specific endpoints that can be used for Discord web hooks for example).

**Note** - This is only tested on Linux. However, technically it should work for Windows as well.

## Command Line Usage
The below is the command line usage from the help menu for this program.

```bash
Help Options
        -cfg= --cfg -cfg <path> > Path to config file override.
        -l --list > Print out full config.
        -v --version > Print out version and exit.
        -h --help > Display help menu.
```

## Configuration
All configuration is done inside a config file on disk via JSON or through environmental variable overrides. The default config path is `/etc/raw/raw.conf` and may be changed via the `-cfg` command line flag.

Any wipe-specific configuration at the top-level of the configuration is used as the default values for each server. Each server may override these by specifiying the same key => value pairs inside of the server array. 

When the program is ran, but no configuration file is found, it will attempt to create the file (likely requiring root privileges by default if trying to create inside of `/etc/`, however). The below are all JSON settings along with their default values, but with added comments. Remember that JSON does **not** support comments. Therefore, copying the below contents with comments will result in errors. This is why it's recommended to allow the program to create a default configuration file.

```
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

    // Timezone for Cron jobs to run with.
    "timezone": "America/Chicago",

    // Either a single string or a slice/array of strings representing when the server should be wiped and processed via Cron format.
    // I would recommend using a Cron generator (there are many online).
    // With that said, the default value wipes at 3:30 PM every Thursday.
    "cronstr": "30 15 * * 4",

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

    // Whether to merge the top-level additional files array and the server-specific files array.
    "deletefilesmerge": false,

    // Additional deletion of files using an array with the array items including "root" as the directory path starting from /home/container (make sure to include beginning /) and another array of strings for file deletion (wildcards included). Ex: [{"root": "/server/rust", "files": [".log"]}]
    "deletefiles": null,

    // Whether to delete server data/files (includes any files with sv.files in the file name). 
    "deletesv": true,

    // Whether to change the world info.
    "changeworldinfo": false,

    // An array with the map, world size, and world seed. (e.g. [{"map": "Procedural Map", "worldsize": 4000, "worldseed": 9213913}])
    "worldinfo": null,

    // World info pick type (1 = pick the next world, otherwise, pick a random world).
    "worldinfopicktype": 1,

    // Whether to merge world information arrays.
    "worldinfomerge": false,

    // Whether to change the hostname.
    "changehostname": true,

    // The hostname format.
    // Would recommend looking here for a cheatsheet on Golang's format library -> https://gosamples.dev/date-time-format-cheatsheet/ 
    // Replacements include:
    // {seconds_left} = Amount of seconds left until next wipe (only valid for warning messages).
    //
    // {tz_one} = Timezone in TTT format (e.g. MST).
    // {tz_two} = Timezone offset in ±hhmm format (e.g. +0100).
    // {tz_three} = Timezone offset in ±hh format (e.g. +01).
    //
    // {month_str_short} = Month in TTT format (e.g. Jan).
    // {month_str_long} = Full month string (e.g. January).
    //
    // {week_day_str_short} = Week day in TTT format (e.g. Mon).
    // {week_day_str_long} = Full week day string (e.g. Monday).
    //
    // {year_one} = Year in TT format (e.g. 22).
    // {year_two} = Full year (e.g. 2022).
    //
    // {month_one} = Month in TT format (e.g. 01).
    // {month_two} = Month in T format (e.g. 1).
    // {month_three} = Month in _T format (e.g. 1).
    //
    // {day_one} = Day in TT format (e.g. 01).
    // {day_two} = Day in T format (e.g. 1).
    // {day_three} = Day in _T format (e.g. 1).
    //
    // {hour_one} = Hour in TT 12 HR format (e.g. 05).
    // {hour_two} = Hour in T 12 HR format (e.g. 5).
    // {hour_three} = Hour in TT 24 HR format (e.g. 17).
    //
    // {min_one} = Minute in TT format (e.g. 06).
    // {min_two} = Minute in T format (e.g. 6).
    //
    // {sec_one} = Second in TT format (e.g. 07).
    // {sec_two} = Second in T format (e.g. 7).
    //
    // {mark_one} = 12-HR mark as TT format (e.g. PM or AM).
    // {mark_two} = 12-HR mark as tt format (e.g. pm or am).
    "hostname": "Vanilla | FULL WIPE {month_two}/{day_two}",
    
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

    // Hooks (submits POST data in JSON with "id", "uuid", "identifier", "name", "ip", and "port" to the hook string/URL). The auth string for each is set to the "Authorization" header (e.g. "Bearer xxxxx").
    // Pre wipe hook.
    prehook: "",
    prehookauth

    // Post wipe hook.
    posthook: "",
    posthookauth: "",

    // Server list (null by default).
   "servers": null
}
```

The servers array includes the following:

```
{
    "servers": [
        {
            // Whether to enable the server or not (enabled by default).
            "enabled": true,

            // The number ID of the server (should be retrieved automatically if the server exists).
            "id": 0,

            // The (short) identifier of the server. Characters before the first "-" in the long UUID.
            "uuid": "",

            // Overrides (retrieve definition from top-level comments above).
            "apiurl": "",
            "apitoken": "",
            "debuglevel": 1,
            "pathtoserverfiles": "/server/rust",
            "timezone": "America/Chicago",
            "cronstr": "30 15 * * 4",
            "deletemap": true,
            "deletebp": true,
            "deletedeaths": true,
            "deletestates": true,
            "deleteidentities": true,
            "deletetokens": true,
            "deletefilesmerge": false,
            "deletefiles": null,
            "deletesv": true,
            "changeworldinfo": false,
            "worldinfo": null,
            "worldinfopicktype": 1,
            "worldinfomerge": false,
            "changehostname": true,
            "hostname": "Vanilla | FULL WIPE {month_two}/{day_two}",
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

With that said, the file `data/egg-rustraw.json` file in this respository is an exported egg that has the standard Rust egg settings from Pterodactyl along with the additional RAW settings. If you import this egg into Pterodactyl under Rust and change the server's egg in the administrator settings, it will easily migrate all existing variable settings over including the Rust egg's default variables.

Each variable is also parsed as a string within the application's code. For true/false, use the 1/0 integers respectfully.

The following is a list of environmental names you can create variables within Pterodactyl Nests/Eggs for overrides.
* **RAW_ENABLED** - Enabled override.
* **RAW_PATHTOFILES** - Path to files override.
* **RAW_TIMEZONE** - Timezone override.
* **RAW_CRONMERGE** - Cron merge override.
* **RAW_CRONSTR** - Cron string override (this is a special case, cron string(s) can be a single string or a string array as a JSON string (e.g. `["*/5 * * * *", "*/2 * * * *"]`)).
* **RAW_DELETEMAP** - Delete map override.
* **RAW_DELETEBP** - Delete blueprints override.
* **RAW_DELETEDEATHS** - Delete deaths override.
* **RAW_DELETESTATES** - Delete states override.
* **RAW_DELETEIDENTITIES** - Delete identities override.
* **RAW_DELETEFILESMERGE** - Whether to merge additional deletion of files override (top-level config and server-specific).
* **RAW_DELETEFILES** - Additional files to delete override. (e.g. `[{"root": "/server/rust", "files": [".log"]}]` as a string).
* **RAW_DELETESV** - Delete server files/data override.
* **RAW_CHANGEWORLDINFO** - Change world info override.
* **RAW_WORLDINFO** - World info override. An array with the map, world size, and world seed. (e.g. `[{"map": "Procedural Map", "worldsize": 4000, "worldseed": 9213913}]` as a string).
* **RAW_WORLDINFOPICKTYPE** - Change world info pick type override.
* **RAW_WORLDINFOMERGE** - Change world info merge override.
* **RAW_CHANGEHOSTNAME** - Change hostname override.
* **RAW_HOSTNAME** - Hostname override.
* **RAW_MERGEWARNINGS** - Merge warnings override.
* **RAW_WARNINGMESSAGES** - Warning messages override (another special case, this should be a JSON string of the normal `warningmessages` JSON item). (e.g. `[{"warningtime": 5, "message": "{seconds_left} until wipe!"}]` as a string).
* **RAW_WIPEFIRST** - Wipe first override.

## Building And Running Project (Manually)
Building the project is simple. We only require `git` and Go.

```bash
# Clone repository.
git clone https://github.com/gamemann/Rust-Auto-Wipe.git

# Change directory to repository.
cd Rust-Auto-Wipe/

# Build using Go into `raw` executable. This will also retrieve needed packages such as Cronv3.
go build -o raw
```

Using `Makefile` as described below is my personal recommendation on installing this application for Linux-based systems.

## Using Makefile + Systemd (Recommended)
You may also use a Makefile I made to build the application and install a Systemd file.

```bash
# Build project (go build -o raw).
make

# Install `rawapp` (/usr/bin/rawapp) and Systemd process.
sudo make install
```

To have the application run on startup, and/or in the background, you may do the following.

```bash
# Reload Systemd daemon (may be needed after install of Systemd service).
sudo systemctl daemon-reload

# Enable (on startup) and start service.
sudo systemctl enable --now raw

# Start service.
sudo systemctl start raw

# Restart service.
sudo systemctl restart raw

# Stop service.
sudo systemctl stop raw

# Disable (on startup) and stop service.
sudo systemctl disable --now raw
```

## Credits
* [Christian Deacon](https://github.com/gamemann)