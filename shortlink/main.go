// Define an App struct
// Contains routers, Middlewares
//
package main

func main() {
	app := App{}
	app.Initialize(getEnv())
	app.Run(":9000")
}
