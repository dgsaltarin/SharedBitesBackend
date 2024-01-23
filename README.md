# SharedBites API

API for sharedbites application

## Deployment

This application run using docker, so in order to run the application use:

```bash
  docker build -t sharedbitesapi .
```

To run the container use:

```bash
  docker run -d -p 8080:8080 --enf-file .env --name sharedbitesapi sharedbitesapi .
```

## API Reference

#### healthcheck

```http
  GET /api/healthcheck
```

Validate that the service is up.

#### Login

```http
  POST /api/items/${id}
```

| Parameter  | Type     | Description                   |
| :--------- | :------- | :---------------------------- |
| `username` | `string` | **Required**. username        |
| `password` | `string` | **Required**. user's password |

Validate user and password and retun a cookie with the authorization token.

#### SignUp

```http
  POST /api/items/${id}
```

| Parameter  | Type     | Description                   |
| :--------- | :------- | :---------------------------- |
| `username` | `string` | **Required**. username        |
| `password` | `string` | **Required**. user's password |
| `email`    | `string` | **Required**. user's email    |

Create a new user in database.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`AWS_ACCESS_KEY` = aws access key

`AWS_SECRET_ACCESS_KEY` = aws secret access key

`SECRET_KEY` = string to sign and validate authorization token
