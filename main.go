package main

func main() {
	path, _ := GetDir("/tests/db.db")
	db, _ := NewDB(path)
	defer db.Close()
	server := NewServer(db)
	server.StartServer()
}
