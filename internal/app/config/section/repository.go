package section

import "github.com/SuperB1aze/catalog-service/internal/app/util"

type (
	Repository struct {
		Postgres RepositoryPostgres
	}

	RepositoryPostgres struct {
		Address      string        `required:"true" split_words:"true"`
		Username     string        `required:"true" split_words:"true"`
		Password     string        `required:"true" split_words:"true"`
		Name         string        `required:"true" split_words:"true"`
		ReadTimeout  util.Duration `default:"5s" split_words:"true"`
		WriteTimeout util.Duration `default:"5s" split_words:"true"`
	}
)
