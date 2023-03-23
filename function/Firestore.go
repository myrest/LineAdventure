package function

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

var (
	fsclient *firestore.Client
	once     sync.Once
)

var firebaseendpoint = ""

func IsRunningInTestMode() bool {
	return firebaseendpoint != ""
}

// GetClient uses sync.Once to ensure that the firestore client is only initialized once
func GetClient() *firestore.Client {
	once.Do(func() {
		// Initialize the client
		//firebaseCredentialFile := Config.ServiceConfig.Firebase.CredentialFile
		firebaseProjectID := "Project-ChatGPT"
		ctx := context.Background()

		//optCredential := option.WithCredentialsFile(firebaseCredentialFile)
		var err error

		if firebaseendpoint != "" {
			os.Setenv("FIRESTORE_EMULATOR_HOST", firebaseendpoint)
		}
		fsclient, err = firestore.NewClient(ctx, firebaseProjectID)
		if err != nil {
			panic(err)
		}
	})
	return fsclient
}

// UpdateOrCreate updates a document with the given id in the specified collection, or creates a new one if it doesn't exist
func UpdateOrCreate(collection string, id string, data interface{}, flag ...bool) (err error) {
	c := GetClient()

	m, err := ConvertStructToMap(data)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		if len(flag) > 0 && flag[0] {
			_, err = c.Collection(collection).Doc(id).Set(context.Background(), m, firestore.MergeAll) //完整覆寫
		} else {
			_, err = c.Collection(collection).Doc(id).Set(context.Background(), m) //合併文件
		}
	}
	return
}

// Delete  deletes a document with the given id in the specified collection
func Delete(collection string, id string) error {
	c := GetClient()
	_, err := c.Collection(collection).Doc(id).Delete(context.Background())
	return err
}

// GetAll returns all documents in the specified collection
func GetAll(collection string) ([]map[string]interface{}, error) {
	c := GetClient()
	iter := c.Collection(collection).Documents(context.Background())
	var result []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result = append(result, doc.Data())
	}
	return result, nil
}

func GetByKeyValue(collection string, key string, value string) ([]map[string]interface{}, error) {
	c := GetClient()
	iter := c.Collection(collection).Where(key, "==", value).Documents(context.Background())
	var result []map[string]interface{}
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		result = append(result, doc.Data())
	}
	return result, nil
}

func GetByID(collection string, ID string) (map[string]interface{}, error) {
	c := GetClient()
	doc, err := c.Collection(collection).Doc(ID).Get(context.Background())
	if err != nil {
		return nil, err
	}
	return doc.Data(), nil
}

func ConvertStructToMap(i interface{}) (map[string]interface{}, error) {
	byt, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(byt, &result); err != nil {
		return nil, err
	}

	filterMap(result)

	return result, nil
}

func filterMap(m map[string]interface{}) {
	for k, v := range m {
		switch v := v.(type) {
		case nil:
			delete(m, k)
		case float64:
			if v == 0 {
				delete(m, k)
			}
		case []interface{}:
			newSlice := make([]interface{}, 0)
			for i := range v {
				if v[i] != nil {
					if innerMap, ok := v[i].(map[string]interface{}); ok {
						filterMap(innerMap)
					}
					newSlice = append(newSlice, v[i])
				}
			}
			if len(newSlice) == 0 {
				delete(m, k)
			} else {
				m[k] = newSlice
			}
		case string:
			if v == "" {
				delete(m, k)
			}
		case map[string]interface{}:
			filterMap(v)
			if len(v) == 0 {
				delete(m, k)
			}
		}
	}
}
