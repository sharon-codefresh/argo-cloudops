package db

import (
	"context"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

type ProjectEntry struct {
	ProjectID  string `db:"project"`
	Repository string `db:"repository"`
}

// Client allows for db crud operations
type Client interface {
	CreateProjectEntry(ctx context.Context, pe ProjectEntry) error
	ReadProjectEntry(ctx context.Context, project string) (ProjectEntry, error)
	DeleteProjectEntry(ctx context.Context, project string) error
}

// SQLClient allows for db crud operations using postgres db
type SQLClient struct {
	host     string
	database string
	user     string
	password string
}

const ProjectEntryDB = "projects"

func NewSQLClient(host, database, user, password string) (SQLClient, error) {
	return SQLClient{
		host:     host,
		database: database,
		user:     user,
		password: password,
	}, nil
}

func (d SQLClient) createSession() (db.Session, error) {
	settings := postgresql.ConnectionURL{
		Host:     d.host,
		Database: d.database,
		User:     d.user,
		Password: d.password,
	}

	return postgresql.Open(settings)
}

func (d SQLClient) CreateProjectEntry(ctx context.Context, pe ProjectEntry) error {
	sess, err := d.createSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	return sess.WithContext(ctx).Tx(func(sess db.Session) error {
		if err := sess.Collection(ProjectEntryDB).Find("project", pe.ProjectID).Delete(); err != nil {
			return err
		}

		if _, err = sess.Collection(ProjectEntryDB).Insert(pe); err != nil {
			return err
		}

		return nil
	})
}

func (d SQLClient) ReadProjectEntry(ctx context.Context, project string) (ProjectEntry, error) {
	res := ProjectEntry{}

	sess, err := d.createSession()
	if err != nil {
		return res, err
	}
	defer sess.Close()

	err = sess.WithContext(ctx).Collection(ProjectEntryDB).Find("project", project).One(&res)
	return res, err
}

func (d SQLClient) DeleteProjectEntry(ctx context.Context, project string) error {
	sess, err := d.createSession()
	if err != nil {
		return err
	}
	defer sess.Close()

	return sess.WithContext(ctx).Collection(ProjectEntryDB).Find("project", project).Delete()
}
