# zerogame
A decentralized cross-platform video game distribution system.

# Instructions

## How to publish a game

### Step 1 - Create an archive
Before publishing your game you must compress it into a .zip archive. One of the files in
the archive must be `install.json` which specifies how to install the archive on various
platforms.

Example:

```
{
  "platforms": [
    {
      "name": "windows",
      "install": ["echo", "installing on windows"],
      "uninstall": ["echo", "uninstalling on widows"],
      "run": ["echo", "running on windows"]
    },
    {  
      "name": "linux",
      "install": ["echo", "installing on linux"],
      "uninstall": ["echo", "uninstalling on linux"],
      "run": ["echo", "running on linux"]
    }
  ]
}
```
Acceptable platform names include:

* windows
* linux
* darwin

### Step 2 - Publish the archive

Your archive must be somewhere publicly accessible on the web: In a shared Dropbox or
Google Drive folder, to give a few examples.

### Step 3 - Generate a feed for your archive

Zerogame distributes software using Feeds, which are JSON descriptions of your archive.
Once your archive is published to the web, you can generate a new feed for your game
using

```
$ zerogame feed
```

Here's an example using Dropbox as the archive storage location:

```
$ zerogame feed
Enter the feed name: my_game
Enter the feed version: 1.0
Enter the feed archive URL: https://www.dropbox.com/s/awg98awe9g7/mygame.zip?dl=1        
Enter the GPG signature URL (optional): https://www.dropbox.com/s/awg98awe9g7/mygame.zip.sig?dl=1
Feed was written to feed.json!
```

### Step 4 - Publish the feed to the web

Publish `feed.json` to the web and share a URL to it. You can optionally store it
alongside your archive.

## How to install a game

Assuming your `feed.json` is publicly available on the web:

```
$ zerogame install https://www.dropbox.com/s/7g707ggaweg/feed.json?dl=1
```

## Troubleshooting

TODO

TODO: Include instructions for signing the archive with gpg and adding to your keyring.