# SuperTokens + Go

A demo implementation of [SuperTokens](https://supertokens.com/) with Golang's [http package](https://pkg.go.dev/net/http).

## General Info

This project aims to demonstrate how to integrate SuperTokens into a Golang server. Its primary purpose is to serve as an educational tool, but it can also be used as a starting point for your own project.

## Repo Structure

### Source

```
📦backend
┣ 📜config.go --> SuperTokens configuration file
┗ 📜main.go --> Entry point of the app
```

#### SuperTokens

The full configuration needed for the SuperTokens' back-end to work is in the `src/config.go` file. This file will differ based on the [auth recipe](https://supertokens.com/docs/guides) you choose.

If you choose to use this as a starting point for your own project, you can further customize SuperTokens in the `src/config.go` file. Refer to our [docs](https://supertokens.com/docs) (and make sure to choose the correct recipe) for more details.

## Application Flow

When using SuperTokens, the front-end never calls directly to the SuperTokens Core, the service that creates and manages sessions. Instead, the front-end calls to the back-end and the back-end calls to the Core. You can read more about [how SuperTokens works here](https://supertokens.com/docs/thirdpartyemailpassword/architecture).

The back-end has two main files:

1. **Entry Point (`main.go`)**

    - Initializes SuperTokens
    - Adds CORS headers for sessions with the front-end
    - Adds SuperTokens middleware
    - Endpoints:
        - `/hello`: Public route not protected by SuperTokens
        - `/sessioninfo`: Uses SuperTokens middleware to pull the session token off the request and get the user session info
        - `/tenants`: Grabs a list of tenants for multitenant configured applications

2. **Configuration (`config.go`)**
    - `supertokensConfig`:        - `supertokens`:
            - `connectionURI`: Now set to `http://localhost:3567` to connect to the self-hosted Docker SuperTokens core. A `docker-compose.yml` file is included to easily start the SuperTokens core and PostgreSQL database locally.
        - `appInfo`: Holds informaiton like your project name
            - `apiDomain`: Sets the domain your back-end API is on. SuperTokens automatically listens to create requests at `${apiDomain}/auth`
            - `websiteDomain`: Sets the domain your front-end website is on
        - `recipeList`: An array of recipes for adding supertokens features

## Self-hosted SuperTokens Core Setup

This project now uses a self-hosted SuperTokens core running in Docker instead of the remote playground instance. To run the application:

1. Start the SuperTokens core and PostgreSQL database:

```bash
docker-compose up -d
```

2. Run the backend server:

```bash
go run .
```

The server will start on `http://localhost:3001` and connect to the SuperTokens core running locally on `http://localhost:3567`.

### Docker Compose Configuration

The `docker-compose.yml` file includes:

1. SuperTokens Core - The authentication service
2. PostgreSQL - The database for storing user information

Both services are configured to restart automatically and include health checks.

## Additional resources

-   Custom UI Example: https://github.com/supertokens/supertokens-web-js/tree/master/examples/react/with-thirdpartyemailpassword
-   Custom UI Blog post: https://supertokens.medium.com/adding-social-login-to-your-website-with-supertokens-custom-ui-only-5fa4d7ab6402
-   Awesome SuperTokens: https://github.com/kohasummons/awesome-supertokens

## Contributing

Please refer to the [CONTRIBUTING.md](https://github.com/supertokens/create-supertokens-app/blob/master/CONTRIBUTING.md) file in the root of the [`create-supertokens-app`](https://github.com/supertokens/create-supertokens-app) repo.

## Contact us

For any questions, or support requests, please email us at team@supertokens.com, or join our [Discord](https://supertokens.com/discord) server.

## Authors

Created with :heart: by the folks at SuperTokens.com.
