package tpl

type Data struct {
	Meta   Meta
	Data   any
	Errors []error
}

type Meta struct {
	Title       string
	Description string
}
