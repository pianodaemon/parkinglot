package dal

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Sets a mongo connection up
func SetUpConnMongoDB(mcli **mongo.Client, uri string) error {

	var err error
	var cli *mongo.Client

	cli, err = mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		goto culminate
	}

	{
		ctxConn, cancelConn := context.WithTimeout(context.Background(),
			10*time.Second)
		defer cancelConn()

		if err = cli.Connect(ctxConn); err != nil {
			goto culminate
		}

		ctx, cancelPing := context.WithTimeout(context.Background(),
			2*time.Second)
		defer cancelPing()

		if err = cli.Ping(ctx, readpref.Primary()); err == nil {
			*mcli = cli
		}
	}

culminate:

	return err
}
