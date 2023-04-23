package main

func main() {
	path, _ := GetDir("/tests")
	idx := NewIndex("cosine")
	db, _ := NewDB(path, idx)
	defer db.Close()
	server := NewServer(db)
	server.StartServer()
}
