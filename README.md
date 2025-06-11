## Planning the REST API 
REST stands for Representational State Transfer





**Planning the API**
A Go-Powered Event Booking REST API 
- `GET/events`
- `GET/events/id`
- `POST/events`
- `PUT/events/id`
-  `DELETE/events/id`
- `POST/signup`
- `POST/login`
- `POST/events/id/register`
- `DELETE/events/id/register`
It will support the following endpoints


## Installing the `gin` framework 
For our project, we are going to use `gin` it is a lightweight middleware that supports `http` requests

To add `gin` framework to our project, we are going to run the following command inside our project directory 

```
go get -u github.com/gin-gonic/gin
```

To use the package inside our project, we can import the gin package as follows

```Go 
import(
	"github.com/gin-gonic/gin"
)
```

and inside our `main` function we can add the line 

```go
server:=gin.Default()
server.Run(":8080")
```

This returns an engine instance with logger and recovery middleware attached which runs on the `localhost:8080`  

## Setting up our first route and handling first request 

There are multiple kind of requests that we have in `http` 
- `GET`
- `POST`
- `DELETE`
- `UPDATE`
In `gin` we can define these requests on the server instance 

```Go
server.GET("/events",getEvents)
```

Here we have given two arguments 
- The endpoint `/events`
- The function to execute `getEvents` whenever this request is placed 

Now we can create our `getEvents` function. Each request in `gin` requires a context, thus while creating any such methods, we are going to pass in the context.

```go
func getEvents(cntxt *gin.Context) {
	context.JSON(http.StatusOK,gin.H{"message":"GET Request"})
}
```

Now when we run our application, we will get the message printed 

## Setting up an `Event` Model

We will create a new package `models` inside this package we will create a new file `event.go` 

Our `Event` struct will have the following attributes
- `ID`
- `Name`
- `Description`
- `Location`
- `DateTime`
- `CreatedBy`
We will also create a slice `events` that will store the `Event` and define two methods that will be used for our events 
- `Save()` - To save the event the slice 
- `GetAllEvents()` - To get all the events

Now we are going to write a method to `POST` an `Event`, for this we are going to create a new handler 

```GO
server.POST("/events",createEvent)
```

and the `createEvent` method is going to be as follows

```GO
func createEvent(cntxt *gin.Context){
	var event models.Event
	err:=cntxt.ShouldBindJSON(&event)
	if err!=nil{
		cntxt.JSON(http.StatusBadRequest, gin.H{"message":"Bad Request","error":err})
	}
	event.ID=1
	event.createdBy=1
	event.Save()
	cntxt.JSON(http.StatusCreated,gin.H{"message":"POST Request","event":event})
}
```

This way we have created the method to `POST` 
## Testing our `POST` requests

To test our `POST` requests and `GET` requests, we are going to create test-files with the extension `http` 

```http 
POST http://localhost:8080/events
content-type:application/json

{
	"name":"Event 1",
	"description":"Description 1",
	"location":"Location 1",
	"dateTime":"2025-01-01T15:30:00.000Z"
}
```

## Initialising the database 
Currently in our system, we are running on a memory based database, that is working for as long as the server is running, however when the server stops running, in that case we want to save the data from the previous session
To solve this we are going to use a database, and we are going to use a SQL database for this

To use the database in our Go project, we need to install the following packages in our project 

```Go
go get github.com/mattn/go-sqlite3
```

 Once we have installed this package, we are going to need to initialise our database, for this we are going to create a new package `db` and inside this package we are going to create a new file called `db.go`
We will import the following packages
- `sql`
- `sqlite3`
We will directly be working with the `sql` package and not the `sqlite3` package, however we don't want the `import` to be removed on compilation, hence while importing we are going to add an `_` in front of the package

We are going to create a pointer to the database instance 

```GO
var DB *sql.DB
```

and we are going to create a function `InitDB()`that will initialise the database for us

```Go
var err error 
DB,err=sql.Open("sqlite3","api.db")
```

Here the `sql.Open()` asks for two things 
- `driverName` - `sqlite3`
- `dataSourceName`-`api.db` where the data will stored at 

If there is an `err ` then we are going to send a `panic` and terminate, else we will set the maximum open connections and idle connections 

In our `main.go` file we are going to add the line 

```Go
db.InitDB()
```
So that as soon as server starts, the database is initialised as well 
## Creating `tables` in our database 

Now that we have initialised our database, we are going to start to create tables in our database, for this we are going to create a new function `createTables`, this function will contain all the logic to keep the tables of the database intact 

```Go 
func createTables(){
 ...
}
```

Since we are using `SQL` database, we are going to be using `SQL` queries quite a lot, for this we are going to write the query to create a table with the same parameters as our `Event` model, we will use this string to execute the query

```SQL
CREATE TABLE IF NOT EXISTS events (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	location TEXT NOT NULL,
	dateTime TEXT NOT NULL,
	user_id INTEGER
)
```

and we are going to execute this query, we will get a result and error on executing, since we are just creating a table, we don't need any result, so we are just going to use the error 

```GO
_,err:=DB.Exec(createEventsTable)
if err!=nil {
	panic("Could not create table"+err.Error())
}
```

This way we have created our table 

## Storing data in our database 
Now that we have setup our tables, we want to store information in the tables, for this we are going to make changes in our `Save()` function 

We are going to first write the query, that we want to use to insert values into our database, since we do not have the information we are going to use `?` and the `gin` package will take care of the rest

```Go
query:=`
	INSERT INTO events(name,description,location,dateTime,user_id) VALUES(?,?,?,?,?)
`
```

Now we are going to prepare this query, with the data we receive from our request 

```Go
stmnt,err:=db.DB.Prepare(query)
if err!=nil{
	return err
}
```

Since we are returning this error, we want the return type of `Save()` to be of `error`
next we are going to execute this statement

```Go 
res,err:=stmnt.Exec(e.Name,e.Description...)
if err!=nil {
	return err
}
id,err:=res.LastInsertID
if err!=nil {
	return err
}
e.ID = id 
return err
```

This will now store the `POST` request data in the table 

Now in our `main.go` file, inside `createEvent`method we are going to make some changes. We will remove the `event.ID=1` as now that will be managed by `AUTOINCREMENT`

```go
err=event.Save()
if err!=nil{
	cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not proceed with request", "error": err.Error()})
}
```

Similarly to display all the rows present in the database, we are going to make changes in the `GetAllEvents()` function in our `Event` package 

```GO
query:=`SELECT * FROM events`
```
We are going to run this query to get our rows, now since are simply querying the database and not updating or inserting any new information we can simply 

```Go
rows,err:=db.DB.Query(query)
```

Since we may have an error, we are going to change the return type of `GetAllEvents`function to `[]Event, error`

We will now create a new `events`slice and use the rows from the slice to fill up the `Event`struct 

```Go
for rows.next() {
	var event Event
	err:=rows.Scan(&event.ID,&event.Name,&event.Description)....
	if err!=nil {
		return nil,err
	}
	events=append(events,event)
}
return events,nil
```

and in our `main.go` file 
```Go 
func getEvents(cntxt *gin.Context){
	events,err:=models.GetAllEvents()
	if err!=nil {
		cntxt.JSON(http.StatusInternalSeverError,gin.H{})
	}
	cntxt.JSON(http.StatusOK,gin.H{"events":events})
}
```

## Deleting the entire data in our database 
To delete the entire records in the database, we are going to write the following query

```sql
DELETE FROM events
```

Now another thing to note is that we still have to make sure that once everything is deleted, we want to start the auto-increment counter from `1` 
For this we are going to write another query which is 

```sql 
DELETE FROM sqlite_sequence WHERE name='events
```

Now that we have both the queries ready, we can execute this query, by creating a new function in `event.go` file named `DeleteAllEvents`

```go
_,err:=db.DB.Exec(query)
if err!=nil {
	return err
}
```

Once we have done this for both the queries, we are going to create a new route and a handler 

```GO
server.DELETE("/events",deleteAllEvents)
```

Once we have this done, we can write the method 

```go
func deleteAllEvents(cntxt *gin.Context) {
	err := models.DeleteAllEvents()
	if err != nil {
	cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "DELETE Request"})
}
```

## Getting a single event by `id` 
Now we are going to write a method to `GET` a single `event` using it's `id`, so we are going to set up a new `route` 

```Go
server.GET("/events/:event_id",getEventByID)
```

Now we are going to set up our method, inside our method, we want to to get the `id` of the record that we want to display, for this we are going to get the parameters from the request itself, which will be 

```Go
cntxt.Param("event_id")
```

Now this will return it, in the form a string, however we need it in the form of `int64` as it is in our `Event`model , so we are going to use `strconv`

```Go
id,err:=strconv.ParseInt(cntxt.Param("event_id"),10,64)
```
 If there is an error we are going to return the `error` as `BadRequest` 

Now that we have our ID, we are going to create a new method in our `event.go` file to fetch the details of an `Event` by it's `id`

For this we are going to create a new function `GetEventByID`

inside this function, we are going to send the `id`, the query to get the record is going to be 

```SQL
SELECT * FROM events WHERE id = ?
```

and since we are returning only a single record we are going to use `QueryRow` method 

```Go
row:=db.DB.QueryRow(query,id)
var event Event
err:=row.Scan(&event.ID,&event.Name,&event.Location...)
if err!=nil {
	return nil,err
}
return &event,nil
```

and once we have done this, inside the `getEventByID` method we are going to 

```GO
event,err!=models.GetEventByID
if err!=nil {
	cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch the event", "error": err.Error()})
}
cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "event": event})
```

## Refactoring our code 
Now we are going to make the code clean and transfer all the logic to `routes` package, for this we are going to create a new package `routes` where we will create a new file `routes.go` inside this we are going to create a function `GetRoutes` which will 
- An instance of the `server` 
- We will define all the routes inside this function 
- We will define all the methods in this file 

## `UPDATE` event 

Now we are going to write the code to `UPDATE` the event, for this we are going to use the `PUT` operation present in `REST` 

```go
server.PUT("/events/:event_id",updateEvent)
```

Once we have created the method, we are again going to follow the same procedure of parsing the value of the `id` from the `request`

So similar to `getEventByID` we are going to fetch the `id` of the event from the `request` itself 

```go
id, err := strconv.ParseInt(cntxt.Param("event_id"), 10, 64)
```

Once we have that `id` we are going to check if there exists an `event` with that `id`

```go
_, err = models.GetEventByID(id)
```

Once this is clear we are going to create a new `event` named `updatedEvent` , where just like creating the event with we are going to bind the event details 

```go 
err = cntxt.ShouldBindJSON(&updatedEvent)
```

and set the `id` of the `updatedEvent` to `id`

```go
 updatedEvent.ID = id
```

Once we have done this, we are going to use the `UpdateEvent` function which we are going to define in `event.go` file 

The query to update the event is going to be 

```sql
UPDATE events
SET name = ?, description=?, location=?, dateTime=?
WHERE id = ?
```
Once we have the query, we are going to follow the exact method we followed while creating an event, that is by using the `Prepare` function 


```go 
stmnt,err:=db.DB.Prepare(query)
if err!=nil {
	return err
}
defer stmnt.Close()
_,err=stmnt.Exec(e.Name,e.Description,e.Location,e.DateTime,e.ID)
if err!=nil {
	return err
}
return nil
```

Once we are done with this function, inside the method, we are going to 

```go 
err=updatedEvent.UpdateEvent()
if err!=nil {
	return internalservererror
}
return success
```

This is how we update the events 

## Deleting `events`
Now we are going to delete event based on `id` , to this first we are going to create a route 


```go 
server.DELETE("/events/:event_id".deleteEventByID)
```

Once we have that done, we are going to create the method, since we are deleting the event by it's `id` we are going to first need the `event` itself, for this we are going to use the same logic we previously used to get the `event`


Once we have set that up, we can go in the `models` file and write the method to delete the event from the database, for this we are going to create a new method `DeleteByID`

In this method, we are going to first write the query which will be 

```SQL
DELETE FROM events WHERE id = ?
```

Now that we have written our query, we are going to prepare the query, it will result in a `stmnt` and `err` if there is any error, we are going to return the `error`

Once that is done, we can execute the statement and move forward. The code for the `deleteByID` is as follows

```Go
func (event Event) DeleteByID() error {
	query := `DELETE FROM events WHERE id=?`
	stmnt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmnt.Close()
	_, err = stmnt.Exec(event.ID)
	if err != nil {
		return err
	}
	return nil
}
```

Now in the `routes` method, it is going to return an `error` 
- If there is an `error` then send the status code `InternalServerError` 
- otherwise return the `StatusOK` 

## Adding a `users` table to the project

Now that we have covered most of the functionality surrounding the event structures, we are going to now work on the `user` aspect of the project

For this we are going to create a `user` struct and set up a table in the database that will store the user data 

For this we are going to go in our `db.go` file and before we create the `events` table we are going to create the `users` table 

For this we are going to first write the query to create the `users` table 

```SQL
CREATE TABLE IF NOT EXISTS users(
id INTEGER PRIMARY KEY AUTOINCREMENT,
email TEXT NOT NULL UNIQUE,
password TEXT NOT NULL
)
```

Once we have this, we are going to make changes in the events table now, we are going to add the following in `events` query

```SQL
user_id INTEGER
FOREIGN KEY(user_id) REFERENCES users(id)
```

Now that we have made these changes, we are going to delete the current database that we have
## Working on the sign up logic

In the `models` package, we are going to add a new file called `users.go`, inside this we are going to define a `struct` that will be as follows, we are going to have the following attributes for the user struct 

- `email`
- `password`
- `id`

The `id` will be taken up automatically, `email` and `password` will be parsed from the request. Now next we are going to define the `Save` method in `users.go` file that will store the users to the database 

Once we have created the `user` struct, we are going to write the method for creating the `user` for our application, this method will be similar in functionality to `createEvent` 

The query to create a new user will be as follows

```SQL 
INSERT INTO users(email,password) VALUES(?,?)
```

Once we have done this, we are going to follow the same approach of preparing the query checking for errors and continue, the code for  `Save()` function is as follows


```Go 
func (u User) Save() error {
	query := `INSERT INTO users(email,password) VALUES(?,?)`
	stmnt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmnt.Close()
	res, err := stmnt.Exec(u.Email, u.Password)
	if err != nil {
		return err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = userID
	return nil
}
```

Now we are going to create a `route ` for this, in the `routes.go` file we are going to create 

```GO 
server.POST("/signup",signUp)
```

and then in `users.go` file we will create the `signUp` method, since we are going to be fetching the data from the request, we are going to bind the data first using `ShouldBindJSON` , then call in the `Save()` function 


```Go
func signUp(cntxt *gin.Context) {
	var user models.User
	err := cntxt.ShouldBindJSON(&user)

	if err != nil {
		cntxt.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse data", "error": err})
	}
	err = user.Save()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	cntxt.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": user})
}

```

## Function to get all the users on the system
To get the list of the users in the system, we are going to follow the same approach we did for retrieving all the `events` 

- Route 
```Go
func getAllUsers(cntxt *gin.Context) {
	users, err := models.GetUsers()
	if err != nil {
		cntxt.JSON(http.StatusInternalServerError, gin.H{"message": "There was some error", "error": err.Error()})
		return
	}
	cntxt.JSON(http.StatusOK, gin.H{"message": "GET Request", "users": users})
}

```

`GetUsers`function 

```Go 
func GetUsers() ([]User, error) {
	query := `SELECT * FROM users`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}
```

Do note that we can add the functionality to get user details, delete users, delete users by `id` similar to how we did with the `events` 
## Hashing the password 

Currently in our application we are storing the passwords as plain text, which is not a good practice as this will allow anyone to see the password of any user present in the system. So we are going to hash the password using encryption 

We are going to get the package in our project 

```Go
go get -u golang.org/x/crypto
```


Now we are going to create a new package for our project `utils` inside this we create a new file called `hash.go` 

Inside this we are going to create a method `hashPassword` for this we are going to 

```Go 
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
```

And now in the `user.go` file before Executing the statement, we are going add the line as follows

```Go
hashedPassword,err:=utils.HashPassword(u.Password);
if err!=nil {
	return err
}
res,err:=stmt.Exec(u.Email,hashedPassword)
```

Now we hashed the password in the database 