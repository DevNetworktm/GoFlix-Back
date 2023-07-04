package easypostgres

type PostgreSQLInit struct {
	request string
	p       *PostgreSQL
}

func NewExecInit(p *PostgreSQL, request string) *PostgreSQLInit {
	return &PostgreSQLInit{
		request: request,
		p:       p,
	}
}

func (pi PostgreSQLInit) Exec() error {
	_, err := pi.p.DB.Exec(pi.request)
	if err != nil {
		return err
	}

	return nil
}
