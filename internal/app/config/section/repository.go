package section

import "time"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address      string        `required:"true" split_words:"true"`
		Username     string        `required:"true" split_words:"true"`
		Password     string        `required:"true" split_words:"true"`
		Name         string        `required:"true" split_words:"true"`
		ReadTimeout  time.Duration `default:"5s" split_words:"true"`
		WriteTimeout time.Duration `default:"5s" split_words:"true"`
	}
)
