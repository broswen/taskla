## Taskla

### Functional Specs
- Store tasks and details within specific groups or the default groups
- create groups with title, description
- create tasks in a group with title, description, status (or default group if not specific)

- user can register with username/password
    - user can authenticate and get jwt
    - user can use jwt to access api

- get all groups by user
    - paginate
    - sort by created date
- get all tasks by user
    - paginate
    - sort by created date
- get all tasks by group and user
    - paginate
    - sort by created date
    - filter by status

- create group
- update group
- delete group (and associated tasks)

- create task
- update task
- delete task



### Project Steps
1. create golang with chi router api

1. create docker image from go api
    1. go testing
    1. swagger documentation

1. store data in postgres database (launch with docker-compose for local testing)

1. add jwt authentication,

1. setup docker-compose for swarm

1. use redis in front of db for lazy loading with ttl and write through

### Technical Specs List

#### Authentication
- [x] register with username and password and registration code
- [ ] username must be unique and password must be complex
- [ ] password is saved using bcrypt
- [x] login with username and password
- [x] receive jwt after successful login, 24hr expiration
- [x] api endpoints authenticate calls using jwt claims

#### Groups
- [x] create groups with id, title, description, created date
- [x] list groups by user (jwt claims)
- [ ] paginate groups, page, limit
- [x] update groups title, description
- [x] delete groups by id, (cascade associated tasks)

#### Tasks
- [x] create task with id, group id, title, description, status, created date
- [x] list tasks by user
    - [ ] paginate tasks, page, limit
- [x] list tasks by group (by user)
    - [ ] paginate tasks, page, limit
- [x] update task title, group id, description, status
- [x] delete task by id

### Technical Specs Implementation Details

#### Authentication
- passwords stored with bcrypt
    - https://pkg.go.dev/golang.org/x/crypto/bcrypt
- JWT created and used for api authentication
    - https://github.com/golang-jwt/jwt

#### Application
- data is stored in postgres database
    - https://pkg.go.dev/github.com/lib/pq
- redis is used for cache
    - https://github.com/go-redis/redis
- zerolog is used for application logging
    - https://github.com/rs/zerolog
- chi router is used for endpoint management
    - https://github.com/go-chi/chi
