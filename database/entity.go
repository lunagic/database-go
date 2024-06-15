package database

type EntityInformation struct {
	TableName string
}

type Entity interface {
	EntityInformation() EntityInformation
}
