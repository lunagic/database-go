package database

type Driver interface {
	DSN() string
	Driver() string
	GenerateInsert(entity Entity) (string, map[string]any, error)
	GenerateDelete(entity Entity) (string, map[string]any, error)
	GenerateSave(entity Entity) (string, map[string]any, error)
	GenerateSelect(entity Entity) (Query, error)
	GenerateUpdate(entity Entity) (string, map[string]any, error)
}
