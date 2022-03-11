package mutations

import (
	"context"
	"fmt"
	"reflect"

	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/util"
)

var ctx = context.Background()

func CreateMinioClient(conn database.DBConn, arg interface{}) (client *minio.Client, err error) {
	var sd database.StorageDeploymentRecord

	switch v := arg.(type) {
	case int64:
		sd, err = database.QueryStorageDeploymentRow(conn, "SELECT * FROM storage_deployments WHERE storage_id = $1", v)
		if err != nil { return nil, util.ProcessErr(err) }
	case database.StorageDeploymentRecord:
		sd = v
	case *database.StorageDeploymentRecord:
		sd = *v
	default:
		return nil, util.ProcessErr(fmt.Errorf("Unexpected input type '%v', expected either 'int64' or 'database.StorageDeploymentRecord'.", reflect.TypeOf(arg)))
	}

	client, err = minio.New(sd.Endpoint, &minio.Options{
        Creds:  credentials.NewStaticV4(sd.AccessKey, sd.SecretKey, ""),
        Secure: sd.UseSSL,
    })
	if err != nil { return nil, util.ProcessErr(err) }
	
	return
}

func CreateMinioAdminClient(conn database.DBConn, arg interface{}) (adminClient *madmin.AdminClient, err error) {
	var sd database.StorageDeploymentRecord

	switch v := arg.(type) {
	case int64:
		sd, err = database.QueryStorageDeploymentRow(conn, "SELECT * FROM storage_deployments WHERE storage_id = $1", v)
		if err != nil { return nil, util.ProcessErr(err) }
	case database.StorageDeploymentRecord:
		sd = v
	case *database.StorageDeploymentRecord:
		sd = *v
	default:
		return nil, util.ProcessErr(fmt.Errorf("Unexpected input type '%v', expected either 'int64' or 'database.StorageDeploymentRecord'.", reflect.TypeOf(arg)))
	}

	adminClient, err = madmin.New(sd.Endpoint, sd.AccessKey, sd.SecretKey, sd.UseSSL)
	if err != nil { return nil, util.ProcessErr(err) }

	return
}
