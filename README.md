## Taskla

Taskla is a REST API backend written in Go with using the Chi router.

User passwords are stored after being hashed with bcrypt and users receive a JWT to authenticate API calls after logging in.


### TODO List

#### Authentication
- [x] register with username and password and registration code
- [ ] username must be unique and password must be complex
- [x] password is saved using bcrypt
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

### MISC
- [ ] load settings and db credentials from config file
- [ ] reconnect db on timeout/disconnect
- [ ] setup docker-compose for swarm
