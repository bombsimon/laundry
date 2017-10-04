# Laundry Booking
This is an application created to host laundry bookings for appartment buildings as a digital soultion. The application can be used even outside appartment buildings as well but might require some tweaking.

## Services
This project will contain two services; one back-end HTTP service and one front-end service. Included with this project are simple tools to generate database schema, populate it, generate PIN numbers etcetera.

## Main features
* Book washing slot online
* Notify user when booked slot is released
* Notify user before or when the slot starts

## Docker
The easiest way to run this would be with docker and therefore I've included a Dockerfile and a docker-compose file.

## Installation
```
$ git clone https://github.com/bombsimon/laundry
$ cd laundry
$ docker-compose up -d
```

### Settings
All the settings related to the server should be located in ```config/back-end.yaml```. Since the file will be copied upon building the container, edit this file before you start the container if you would like to make any changes. Settings related to booking and such will be stored in a database and configurable.

# TODO
## First release
* Better slog management
* Determine if I should use something other than MySQL. MongoDB?
* All configuration not related to server/port should be configurable via GUI (stored in DB)
* Authorization
  * JWT + validate in middleware
  * Login? PIN only? RFID?
* Logs and history
* Create tool to generate base data such as machines
* Watch/notification/reminders
  * Remind bookers via mail/SMS about times
  * Notify users watching a specific time
* Frontend

## Future
There is a lot of things I would like to do with this project but as of now I've just put them in the future category. The things I would like to see the most
* Hook to an actual digital lock to unlock doors
* RFID support or similar
* WordPress plugin
* Integrate with digital locks
