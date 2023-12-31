package leveldb_admin

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"net/http"
	"strconv"
)

type keyListRes struct {
	Items      []string
	SearchText string
	IsPart     bool
}

func (l *LevelAdmin) apiKeys(writer http.ResponseWriter, request *http.Request) {
	db := request.URL.Query().Get("db")
	if db == "" {
		http.NotFound(writer, request)
		return
	}

	prefix := request.URL.Query().Get("prefix")
	searchText := request.URL.Query().Get("searchText")
	limitStr := request.URL.Query().Get("limit")
	limit := 15
	if limitStr != "" {
		limitRe, err := strconv.Atoi(limitStr)
		if err != nil {
			l.writeError(writer, err)
		}
		limit = limitRe
	}

	if limit > 15 {
		limit = 15
	}

	if limit < 0 {
		limit = 15
	}

	res := &keyListRes{IsPart: false}

	if load, ok := l.dbs.Load(db); ok {
		db := load.(*leveldb.DB)

		iter := db.NewIterator(util.BytesPrefix(l.keySerializer.Deserialize(prefix)), nil)
		defer iter.Release()

		if searchText != "" {
			iter.Seek(l.keySerializer.Deserialize(searchText))
		}

		for iter.Next() {
			if len(res.Items) >= limit {
				res.SearchText = l.keySerializer.Serialize(iter.Key())
				res.IsPart = true

				l.writeJson(writer, res)
				return
			}

			res.Items = append(res.Items, l.keySerializer.Serialize(iter.Key()))
		}

		err := iter.Error()
		if err != nil {
			l.writeError(writer, err)
			return
		}

		res.IsPart = false
		l.writeJson(writer, res)
	} else {
		http.NotFound(writer, request)
	}
}
