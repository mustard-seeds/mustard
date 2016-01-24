package main
import (
    "fmt"
    "mustard/crawl/handler"
)
func main() {
    handler.InitCrawlService()
    var input string
    fmt.Scanln(&input)
}
