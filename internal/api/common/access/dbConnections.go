package access

import "checkmarks/internal/config"

type DbConnections struct {
	Mongo *Mongo
}

func NewServiceDbConnections(c *config.Config) (*DbConnections, error) {

	mongo, err := initMongoConnection(c.Mongo)

	if err != nil {
		return nil, err
	}

	sdc := DbConnections{
		mongo,
	}

	return &sdc, nil
}
