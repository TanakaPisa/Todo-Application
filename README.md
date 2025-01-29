# Todo-Application
A to-do item should include a description and a status of "not started", "started", or "completed". Users should be able to display a list of all to-do items including their statuses, and should be able to update a to-do item to change its status or delete it entirely. Feel free to add any additional functionality if you want.


1) Basic CLI application
- Create a command line application that uses flags to accept a to-do item adds it to an empty list of to-do items and prints the list to console
- After printing the list of to-do items, save them to a file on disk
- When the application starts, load all to-do items from disk before adding new item
- Allow the user to update the description of a to-do item or delete it
2) Advanced functionality
- Use the "log/slog" structured logging package to log errors and when data is saved to disk
- Use the "context" package to add a TraceID to enable traceability of calls through the solution by adding it to all logs
- Separate the core todo store logic into a different package/module to main/CLI code
- Write unit tests to cover usefully testable code
- Use the "os/signal" package and ensure that the application only exits when it receives the interrupt signal (ctrl+c)
3) API
- Use ServeMux in the "net/http" package to expose json http endpoints: "/create", "/get", "/update", and "/delete"
- Add a middleware to create a context with TraceID
4) Web page
- Use "http.FileServer" to serve a static page to a new "/about" endpoint
- Use "html/template" to serve a dynamic page containing a list of all to-do items to a new "/list" endpoint
5) Concurrency
- Use the Actor/Communicating Sequential Processes (CSP) pattern to support concurrent reads and concurrent safe write
- Use Parallel tests to validate that the solution is concurrent safe


# CLI operations
1) Add new item
- go run main.go -action=add -id=1 -desc="new Item" -status="not started"
2) Update item
- go run main.go -action=update -id=1 -desc="new Item" -status="started"
3) Remove item 
- go run main.go -action=remove -id=0

# API operations
1) Add new item
- curl -X POST -H "Content-Type: application/json" -d '{"id":1, "desc":"Buy bread", "status":"done"}' http://localhost:8080/create
2) Update new item
- curl -X POST -H "Content-Type: application/json" -d '{"id":1, "desc":"Buy bread", "status":"done"}' http://localhost:8080/update
3) remove an item
- curl -X POST -H "Content-Type: application/json" -d '{"id":1}' http://localhost:8080/delete
4) get an item
- curl -X POST -H "Content-Type: application/json" -d '{"id":1}' http://localhost:8080/get