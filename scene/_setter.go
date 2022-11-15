package scene

type Setter interface {
	Default() Setter
	Current() Setter
	Set(Setter)
	Load(any) (Setter, error)
}

// func LoadSettings(fsys fs.FS, name string, setter Setter) {
// 	setter.Set(setter.Default())
// 	setter.Set(data)
// }
