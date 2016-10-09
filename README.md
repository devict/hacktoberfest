# hacktoberfest

A little web app for tracking participation in the devICT Hacktoberfest event.
The entire idea and name "Hacktoberfest" comes from the [Digital Ocean and
GitHub event](https://hacktoberfest.digitalocean.com/) by the same name.

## contribute

Make Wichita a better place through code. Hop on over to [Our Hacktoberfest app](https://devict-hacktoberfest.herokuapp.com) and register to register with your GitHub portfolio.

# Development

You need to create an application on GitHub
[here](https://github.com/settings/applications/new) then put the ID
and secret in a file `secret.env` like this

```
GITHUB_KEY=123abc123abc
GITHUB_SECRET=123abc123abc123abc123abc
```

# Running

You can run the app locally using [Docker](https://docker.com). There is
a docker-compose file that will run the app and database.

Using the configured compose file if you edit anything under `public/`
or `templates/` then the change will be available immediately. If you
change any `.go` files you will have to rebuild the image.

I'd recommend starting the db first in daemon mode then start the web
service in the foreground so you can easily restart it by canceling with
`Ctrl-c` then restart.

    docker-compose up -d db
    docker-compose up --build web
