# Hacktoberfest

A little web app for tracking participation in the devICT Hacktoberfest event.

## Participate

Make Wichita a better place through code. Hop on over to the
[devICT Hacktoberfest app](https://devict-hacktoberfest.herokuapp.com) and
register with your GitHub profile.

# Development

Follow the steps below to get the app running locally on your system.

## Installation

Download & install
[Docker](https://docs.docker.com/install/#supported-platforms) on your system.

Create an application from your
[GitHub Settings page](https://github.com/settings/applications/new). Set the
**Homepage URL** to `http://localhost:8080` and **Authorization callback URL**
to `http://localhost:8080/auth/github`.

Copy the `secret.env.example` file and rename it to `secret.env`:

```
cp secret.env.example secret.env
```

Paste the **Client ID** and **Client Secret** from your registered GitHub
application into `secret.env`.

Start the database in daemon mode:

```
docker-compose up -d db
```

Start the web service:

```
docker-compose up --build web
```

Install frontend assets:

```
npm install
```

Build the frontend assets:

```
npm run production
```

You're ready to go! Visit [http://localhost:8080](localhost:8080) in your
browser.

## Developing Locally

Using the configured Docker Compose file if you edit anything under `public/` or
`templates/` then the change will be available immediately. If you change any
`.go` files you will have to rebuild the image. While the web service is running
press <kbd>ctrl</kbd> + <kbd>c</kbd> to cancel the process and then run
`docker-compose up --build web` again to restart.

Running `npm run watch` will watch for frontend changes and rebuild CSS assets.

This website uses the Typicons icon library. [Download the SVG files from their website](https://www.s-ings.com/typicons/).

## Credit

The entire idea and name "Hacktoberfest" comes from the
[Digital Ocean event](https://hacktoberfest.digitalocean.com) by the same name.
