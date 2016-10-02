# hacktoberfest

A little web app for tracking participation in the devICT Hacktoberfest event.
The entire idea and name "Hacktoberfest" comes from the [Digital Ocean and
GitHub event](https://hacktoberfest.digitalocean.com/) by the same name.

# Development

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
