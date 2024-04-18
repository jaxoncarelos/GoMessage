# GoMessage
An easily expandable Go TCP messaging framework

## How it works
The setup is relatively simple, and easily modifable. By default, it comes with no commands setup. Add commands by doing
```go
server := gomessage.NewGoMessage()
server.AddCommand("aliasHere", func(str string, conn Net.Conn) {
  // logic here
}
// if the message the user sends is setup as aliasHere:argYouWant:anotherArgYouWant it will run the logic in here AFTER the OnMessage.
```
You can add more basic logic that you want to do on every connection or mesasge by assigning a function to the servers onMessage and onConnection field.
```go
server := gomessage.NewGoMessage()
server.OnConnect = func(net.Conn) {
  fmt.Println("New connection!")
}
server.OnMessage = func(str string, conn net.Conn) {
  fmt.Printf("New message from %s: %s", conn.LocalAddr.String(), str)
}
```
This logic will be ran on every connection, and every message.
