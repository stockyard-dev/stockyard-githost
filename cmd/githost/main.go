package main
import ("fmt";"log";"net/http";"os";"github.com/stockyard-dev/stockyard-githost/internal/server";"github.com/stockyard-dev/stockyard-githost/internal/store")
func main(){port:=os.Getenv("PORT");if port==""{port="9710"};dataDir:=os.Getenv("DATA_DIR");if dataDir==""{dataDir="./githost-data"}
db,err:=store.Open(dataDir);if err!=nil{log.Fatalf("githost: %v",err)};defer db.Close();srv:=server.New(db)
fmt.Printf("\n  Githost — Git repository host\n  Dashboard:  http://localhost:%s/ui\n  API:        http://localhost:%s/api\n\n",port,port)
log.Printf("githost: listening on :%s",port);log.Fatal(http.ListenAndServe(":"+port,srv))}
