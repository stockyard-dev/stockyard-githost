package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Repository struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	Visibility string `json:"visibility"`
	SizeBytes int `json:"size_bytes"`
	StarCount int `json:"star_count"`
	Status string `json:"status"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"githost.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS repositories(id TEXT PRIMARY KEY,name TEXT NOT NULL,description TEXT DEFAULT '',default_branch TEXT DEFAULT 'main',visibility TEXT DEFAULT 'private',size_bytes INTEGER DEFAULT 0,star_count INTEGER DEFAULT 0,status TEXT DEFAULT 'active',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Repository)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO repositories(id,name,description,default_branch,visibility,size_bytes,star_count,status,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Description,e.DefaultBranch,e.Visibility,e.SizeBytes,e.StarCount,e.Status,e.CreatedAt);return err}
func(d *DB)Get(id string)*Repository{var e Repository;if d.db.QueryRow(`SELECT id,name,description,default_branch,visibility,size_bytes,star_count,status,created_at FROM repositories WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Description,&e.DefaultBranch,&e.Visibility,&e.SizeBytes,&e.StarCount,&e.Status,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Repository{rows,_:=d.db.Query(`SELECT id,name,description,default_branch,visibility,size_bytes,star_count,status,created_at FROM repositories ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Repository;for rows.Next(){var e Repository;rows.Scan(&e.ID,&e.Name,&e.Description,&e.DefaultBranch,&e.Visibility,&e.SizeBytes,&e.StarCount,&e.Status,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *Repository)error{_,err:=d.db.Exec(`UPDATE repositories SET name=?,description=?,default_branch=?,visibility=?,size_bytes=?,star_count=?,status=? WHERE id=?`,e.Name,e.Description,e.DefaultBranch,e.Visibility,e.SizeBytes,e.StarCount,e.Status,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM repositories WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM repositories`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]Repository{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR description LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,description,default_branch,visibility,size_bytes,star_count,status,created_at FROM repositories WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []Repository;for rows.Next(){var e Repository;rows.Scan(&e.ID,&e.Name,&e.Description,&e.DefaultBranch,&e.Visibility,&e.SizeBytes,&e.StarCount,&e.Status,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM repositories GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
