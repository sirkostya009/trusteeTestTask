# Trustee Golang Test Task

### How to run:

1. Make sure you have a GCP project with Calendar API enabled and a Desktop App OAuth2.0 client ID active.
Download `credentials.json` and drop it in the root of the project.
Make sure the `redirect_uris` array contains `http://localhost:3000/api/authenticate`, for redirecting purposes.
2. After that, do `go run cmd/main.go`

### How to use:

1. Hit `localhost:3000/authenticate` in the browser, authenticate with the Google account you have the Calendar API enabled for.
   1. Copy the `google-session` cookie from the response headers and paste it in the request headers for the next requests.
2. Hit `GET localhost:3000/api/showtime`, which should already show available showtime slots filtered after your
calendar events.
   1. Hit `PATCH localhost:3000/api/showtime/:id` to book the showtime seat, and add the event to your calendar.
   2. Hit `DELETE localhost:3000/api/showtime/:id` to cancel booking, and clear the calendar event.

And that's it!
