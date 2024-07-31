# Visible
A lightweight library to control the visibility of struct fields based on roles.

## Installation
```bash
go get github.com/fattymango/visible
```

## Examples
```go
import "github.com/fattymango/visible"

type User struct {
    ID int `json:"id"`
    Name string `json:"name"`
    AdminData interface{} `json:"admin_data" visible:"admin"`
}



func main() {
   user := User{
    ID: 1,
    Name: "John Doe",
    AdminData: "Admin Data",
}

res,err := visible.CleanStruct(user, "admin")
if err != nil {
    log.Fatal(err)
}
fmt.Println(res)
}
// Output: {id:1 name:John Doe admin_data:Admin Data}

res,err := visible.CleanStruct(user, "user")
if err != nil {
    log.Fatal(err)
}
fmt.Println(res)
}
// Output: {id:1 name:John Doe}

```


## Example with Casbin and GoFiber
```go
import (
	"github.com/fattymango/visible"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

func NewCasbinMiddleware(e *casbin.Enforcer, resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// get role from path
		role := c.Params("role")

		if role == "" {
			return c.SendStatus(fiber.StatusForbidden)
		}

		if ok, _ := e.Enforce(role, resource, action); ok {
			return c.Next()

		}
		return c.SendStatus(fiber.StatusForbidden)

	}
}

type Data struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// AliceData is only visible to alice
	AliceData string `json:"alice_data" visible:"alice"`
	// BobData is only visible to bob
	BobData string `json:"bob_data" visible:"bob"`
	// AliceAndBobData is visible to both alice and bob
	AliceAndBobData string `json:"alice_and_bob_data" visible:"alice,bob"`
	// this field will never be visible to anyone
	PrivateData string `json:"-"`
}

func NewSuccessResponse(ctx *fiber.Ctx, data interface{}) error {

	cleanData, err := visible.CleanStruct(data, ctx.Params("role"))
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	return ctx.Status(200).JSON(cleanData)
}
func main() {
	dsn := "test:test@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	a, err := gormadapter.NewAdapterByDB(db) // Your driver and data source.
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("casbin_model.conf", a)
	if err != nil {
		panic(err)
	}
	e.LoadPolicy()
	e.AddPolicy("alice", "data1", "read")
	e.AddPolicy("bob", "data1", "read")
	e.AddPolicy("charlie", "data1", "read")
	app := fiber.New()

	// examples:
	// http://localhost:3000/alice --> {ID: 1, Name: "alice", AliceData: "alice_data", AliceAndBobData: "alice_and_bob_data"}
	// http://localhost:3000/bob --> {ID: 1, Name: "bob", BobData: "bob_data", AliceAndBobData: "alice_and_bob_data"}
	// http://localhost:3000/charlie --> {ID: 1, Name: "charlie"}
	// http://localhost:3000/saif --> 403 Forbidden
	app.Get("/:role",
		NewCasbinMiddleware(e, "data1", "read"),
		func(ctx *fiber.Ctx) error {
			data := Data{
				ID:              1,
				Name:            c.Params("role"),
				AliceData:       "alice_data",
				BobData:         "bob_data",
				AliceAndBobData: "alice_and_bob_data",
				PrivateData:     "private_data",
			}

			return NewSuccessResponse(ctx, data)
		})

	app.Listen(":3000")
}
```