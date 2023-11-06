<p align="center">
  <a href="https://github.com/polldo/govod">
    <img src="https://private-user-images.githubusercontent.com/17302582/280764178-b413f94c-b050-4a72-ac69-dc1c223df3b7.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTEiLCJleHAiOjE2OTkyOTIzNzgsIm5iZiI6MTY5OTI5MjA3OCwicGF0aCI6Ii8xNzMwMjU4Mi8yODA3NjQxNzgtYjQxM2Y5NGMtYjA1MC00YTcyLWFjNjktZGMxYzIyM2RmM2I3LnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFJV05KWUFYNENTVkVINTNBJTJGMjAyMzExMDYlMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjMxMTA2VDE3MzQzOFomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTA0OGM4YTk3MWIwYjk0YWZjNWMxMGE5NjBmNmI3NjA0YzRkYjQ0NmI2ZDViZmU5NWYwMDE5ODM0NTM3MzU2ZDEmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0JmFjdG9yX2lkPTAma2V5X2lkPTAmcmVwb19pZD0wIn0.wXVek7UfWD8QaJfK-wtMTUtk5pX3WVEOl6CfmvnFdYY" />
  </a>
</p>

<h1 align="center">govod</h1>
<p align="center"><i><b>A simple web platform to sell videos on demand.</b></i></p>
<hr>


## Features
- Login with google or password.
- Require email activation.
- Password reset.
- Free samples.
- Shopping cart.
- Purchase with stripe or paypal.
- Play videos through [VideoJS](https://github.com/videojs) (support all major streaming formats).
- Store video progress.

| <img src="https://github.com/polldo/govod/assets/17302582/97360e90-4924-459b-b9d7-787156bb3b4d" width="300" alt=""/> Free sample demo | <img src="https://github.com/polldo/govod/assets/17302582/57195757-6c48-42bf-b670-f92a5ffbdadd" width="300" alt=""/> Video player demo | <img src="https://github.com/polldo/govod/assets/17302582/2b4eb834-8d07-4b7d-987c-03a5bcb5ac68" width="300" alt=""/> Shopping cart demo |
|-------------------|-------------------|-------------------|

## Configuration
Backend: These are some environment variables you might need to set. You can find them
all in the [configuration
package](https://github.com/polldo/govod/blob/main/config/config.go).
```bash
# Web configuration.
export GOVOD_WEB_ADDRESS="127.0.0.1:8000"
export GOVOD_AUTH_ACTIVATION_REQUIRED=true
# Database configuration.
export GOVOD_DB_USER="postgres"
export GOVOD_DB_NAME="govod"
# SMTP configuration.
export GOVOD_EMAIL_HOST=""
export GOVOD_EMAIL_PORT=""
export GOVOD_EMAIL_ADDRESS=""
export GOVOD_EMAIL_PASSWORD=""
# Paypal configuration.
export GOVOD_PAYPAL_CLIENT_ID=""
export GOVOD_PAYPAL_SECRET=""
# Stripe configuration.
export GOVOD_STRIPE_API_SECRET=""
export GOVOD_STRIPE_WEBHOOK_SECRET=""
# Google oauth configuration.
export GOVOD_OAUTH_GOOGLE_CLIENT=""
export GOVOD_OAUTH_GOOGLE_SECRET=""
export GOVOD_OAUTH_GOOGLE_URL=""
export GOVOD_OAUTH_GOOGLE_REDIRECT_URL=""
export GOVOD_OAUTH_LOGIN_REDIRECT_URL=""
# CORS configuration.
export GOVOD_CORS_ORIGIN="http://mylocal.com:3000"
```

Frontend: set these env vars in a `frontend/.env.local` file:
```bash
NEXT_PUBLIC_PAYPAL_CLIENT_ID=""
NEXT_PUBLIC_BASE_URL="http://mylocal.com:8000"
```


## Build and run
This project uses PostgreSQL, so you need it to run the backend.

You'll also need to run the [db migrations](https://github.com/polldo/govod/tree/main/database/sql/migration).

Then, to correctly integrate stripe and paypal, you need to make an account and
fill the environment variables accordingly.

For the SMTP server you can use a dedicated service like Mailtrap.

To run the backend - assuming your env variables are defined in a .env file:
```bash
cd cmd/server
. ./.env
go run .
```

To run the frontend:
```bash
cd frontend
npm run dev
```

