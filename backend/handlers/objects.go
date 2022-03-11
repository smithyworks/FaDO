package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/smithyworks/FaDO/database"
	"github.com/smithyworks/FaDO/mutations"
	"github.com/smithyworks/FaDO/util"
)

func Objects(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
		if path := r.URL.Query().Get("path"); path != "" {
			conn, err := database.Acquire()
			if err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer conn.Release()

			bucketName := strings.Split(path, "/")[0]
			if len(path) == len(bucketName) {
				util.PrintErr(fmt.Errorf("Expected path value to be of the form <bucket-name>/<object-name>, but got %v.", path))
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			objectName := path[len(bucketName) + 1:]

			if bucket, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE name = $1", bucketName); err != nil {
				util.PrintErr(err)
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			} else {
				if object, err := database.QueryObjectRow(conn, "SELECT * FROM objects WHERE name = $1 AND bucket_id = $2", objectName, bucket.BucketID); err != nil {
					util.PrintErr(err)
					http.Error(w, "Not Found", http.StatusNotFound)
					return
				} else {
					ServeObject(w, r, conn, bucket, object)
					return
				}
			}
		} else {
			SendResources(w)
			return
		}
    } else if r.Method == "POST" {
		r.ParseMultipartForm(100 << 20)

		bucketId, err := strconv.Atoi(r.FormValue("bucket"))
		if err != nil {
			util.PrintErr(err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		var bucket database.BucketRecord
		var client *minio.Client
		if conn, err := database.Acquire(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer conn.Release()

			if b, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE bucket_id = $1", bucketId); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				bucket = b
			}

			if c, err := mutations.CreateMinioClient(conn, bucket.StorageID); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				client = c
			}

			file, handler, err := r.FormFile("file")
			if err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer file.Close()
	
			if _, err = client.PutObject(ctx, bucket.Name, handler.Filename, file, handler.Size, minio.PutObjectOptions{}); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if _, err = database.InsertObject(conn, database.ObjectRecord{BucketID: bucket.BucketID, Name: handler.Filename}); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}
		
		SendResources(w)
		return
	} else {
		util.PrintErr(fmt.Errorf("Method not supported. Got %v.", r.Method))
		http.Error(w, "Method Not Supported", http.StatusNotFound)
		return
	}
}

func Object(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
		object_id, err := strconv.Atoi(mux.Vars(r)["object_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/objects/<int>", r.URL.RequestURI()))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		conn, err := database.Acquire()
		if err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer conn.Release()

		if object, err := database.QueryObjectRow(conn, "SELECT * FROM objects WHERE object_id = $1", object_id); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			if bucket, err := database.QueryBucketRow(conn, "SELECT * FROM buckets WHERE bucket_id = $1", object.BucketID); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				ServeObject(w, r, conn, bucket, object)
				return
			}
		}
    } else if r.Method == "DELETE" {
		object_id, err := strconv.Atoi(mux.Vars(r)["object_id"])
		if err != nil {
			util.PrintErr(fmt.Errorf("Path not found. Expected %v, got %v.", "/api/objects/<int>", r.URL.RequestURI()))
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if tx, err := database.Begin(); err != nil {
			util.PrintErr(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		} else {
			defer tx.Rollback(ctx)

			if object, err := database.QueryObjectRow(tx, "SELECT * FROM objects WHERE object_id = $1", object_id); err != nil {
				util.PrintErr(err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			} else {
				if err = mutations.DeleteObject(tx, object); err != nil {
					util.PrintErr(err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				} else {
					tx.Commit(ctx)
				}
			}
		}

		SendResources(w)
		return
	} else {
		util.PrintErr(fmt.Errorf("Method not supported. Got %v.", r.Method))
		http.Error(w, "Method Not Supported", http.StatusNotFound)
		return
	}
}

func ServeObject(w http.ResponseWriter, r *http.Request, conn database.DBConn, bucket database.BucketRecord, object database.ObjectRecord) {
	var client *minio.Client
	if c, err := mutations.CreateMinioClient(conn, bucket.StorageID); err != nil {
		util.PrintErr(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		client = c
	}

	if minioObj, err := client.GetObject(ctx, bucket.Name, object.Name, minio.GetObjectOptions{}); err != nil {
		util.PrintErr(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	} else {
		defer minioObj.Close()
		log.Printf("%v", minioObj)
		w.Header().Set("Content-Disposition", "attachment; filename=" + strconv.Quote(object.Name))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeContent(w, r, object.Name, time.UnixMicro(0), minioObj)
		return
	}
}
