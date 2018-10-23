package goa

func middleware1(c *Context){
    println("middleware1")
}

func middleware2(c *Context){
    println("middleware2 pre")
    c.Next()
    println("middleware2 post")
}

func middleware3(c *Context){
    c.Next()
    println("middleware3")
}
