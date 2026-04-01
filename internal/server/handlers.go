package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-githost/internal/store")
func(s *Server)handleListRepos(w http.ResponseWriter,r *http.Request){list,_:=s.db.ListRepos();if list==nil{list=[]store.Repo{}};writeJSON(w,200,list)}
func(s *Server)handleCreateRepo(w http.ResponseWriter,r *http.Request){var repo store.Repo;json.NewDecoder(r.Body).Decode(&repo);if repo.Name==""{writeError(w,400,"name required");return};if repo.Visibility==""{repo.Visibility="private"};if repo.DefaultBranch==""{repo.DefaultBranch="main"};s.db.CreateRepo(&repo);writeJSON(w,201,repo)}
func(s *Server)handleDeleteRepo(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.DeleteRepo(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleListIssues(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);list,_:=s.db.ListIssues(id);if list==nil{list=[]store.Issue{}};writeJSON(w,200,list)}
func(s *Server)handleCreateIssue(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var iss store.Issue;json.NewDecoder(r.Body).Decode(&iss);iss.RepoID=id;if iss.Title==""{writeError(w,400,"title required");return};s.db.CreateIssue(&iss);writeJSON(w,201,iss)}
func(s *Server)handleCloseIssue(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.CloseIssue(id);writeJSON(w,200,map[string]string{"status":"closed"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
