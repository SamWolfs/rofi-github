package github

type Resource interface {
	Format(string) string
	View()
}
