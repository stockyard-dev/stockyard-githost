package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Repo struct{
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	DefaultBranch string `json:"default_branch"`
	Visibility string `json:"visibility"`
	CloneURL string `json:"clone_url"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"githost.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS repos(id TEXT PRIMARY KEY,name TEXT NOT NULL,description TEXT DEFAULT '',default_branch TEXT DEFAULT 'main',visibility TEXT DEFAULT 'private',clone_url TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Repo)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO repos(id,name,description,default_branch,visibility,clone_url,created_at)VALUES(?,?,?,?,?,?,?)`,e.ID,e.Name,e.Description,e.DefaultBranch,e.Visibility,e.CloneURL,e.CreatedAt);return err}
func(d *DB)Get(id string)*Repo{var e Repo;if d.db.QueryRow(`SELECT id,name,description,default_branch,visibility,clone_url,created_at FROM repos WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Description,&e.DefaultBranch,&e.Visibility,&e.CloneURL,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Repo{rows,_:=d.db.Query(`SELECT id,name,description,default_branch,visibility,clone_url,created_at FROM repos ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Repo;for rows.Next(){var e Repo;rows.Scan(&e.ID,&e.Name,&e.Description,&e.DefaultBranch,&e.Visibility,&e.CloneURL,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM repos WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM repos`).Scan(&n);return n}
