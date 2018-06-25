# Laundry Booking
This is an application created to host laundry bookings for appartment buildings
as a digital soultion. The main reason for this is to provide a simple *deploy
it yourself* solution as an alternative to the expensive existing services. At
the same time I would like to priovide features that most current systems are
missing surch as remider notification, booking proposals and similar features.

This service is a JSON RESTful API which should be combined with a front-end
service to provide the best user experience. 

## Main features
* Book washing slots online
* Notify user when booked slot is released
* Notify user before or when the slot starts
* Enable reminders to book a new slot

## Docker
The easiest way to run this would be with docker and therefore I've included a
Dockerfile and a docker-compose file.

## Installation
```
$ git clone https://github.com/bombsimon/laundry
$ cd laundry
$ docker-compose up -d
```

### Settings
All the settings related to the server should be located in
`config/back-end.yaml`. Since the file will be copied upon building the
container, edit this file before you start the container if you would like to
make any changes. Settings related to booking and such will be stored in a
database and configurable.

## TODO
### First release
* Better log management
* Determine if I should use something other than MySQL
* All configuration not related to server/port should be configurable via GUI
  (stored in DB)
* Authorization
  * JWT + validate in middleware?
  * PIN only login
* Logs and history
* Create tool to generate base data such as machines
* Watch/notification/reminders
  * Remind bookers via mail/SMS about times
  * Notify users watching a specific time

### Future
There is a lot of things I would like to do with this project but as of now I've
just put them in the future category. The things I would like to see the most
* Hook to an actual digital lock to unlock doors
* RFID support or similar
* WordPress plugin
* Integrate with digital locks
