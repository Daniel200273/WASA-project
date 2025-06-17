package database

// SetName is an example that shows you how to execute insert/update (compatibility method)
func (db *appdbimpl) SetName(name string) error {
	_, err := db.c.Exec("INSERT OR REPLACE INTO example_table (id, name) VALUES (1, ?)", name)
	return err
}
