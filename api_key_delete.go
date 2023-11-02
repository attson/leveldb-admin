package leveldb_admin

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"net/http"
)

type deleteRes struct {
	Success bool
}

type deleteReq struct {
	DB  string
	Key string
}

func (l *LevelAdmin) apiKeyDelete(writer http.ResponseWriter, request *http.Request) {
	reqData := &deleteReq{}
	err := json.NewDecoder(request.Body).Decode(&reqData)
	if err != nil {
		l.writeError(writer, err)

		return
	}

	if reqData.DB == "" || reqData.Key == "" {
		http.NotFound(writer, request)
		return
	}

	if load, ok := l.dbs.Load(reqData.DB); ok {
		db := load.(*leveldb.DB)
		if has, err := db.Has(l.keySerializer.Deserialize(reqData.Key), nil); has && err == nil {
			db.Delete(l.keySerializer.Deserialize(reqData.Key), nil)
		} else {
			http.NotFound(writer, request)
			return
		}

		l.writeJson(writer, &deleteRes{Success: true})
	} else {
		http.NotFound(writer, request)
	}
}
