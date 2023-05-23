package models

import (
	"ChatDanBackend/config"
	"ChatDanBackend/utils"
	"github.com/meilisearch/meilisearch-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
)

var meilisearchClient *meilisearch.Client

func InitSearch() {
	var err error
	if config.Config.MeilisearchUrl == "" {
		return
	}
	meilisearchClient = meilisearch.NewClient(meilisearch.ClientConfig{
		Host:   config.Config.MeilisearchUrl,
		APIKey: config.Config.MeilisearchApiKey,
	})
	utils.Logger.Info("Meilisearch initialized")

	// create or update indexes
	var searchModels = []SearchModel{BoxSearchModel{}, TagSearchModel{}}

	for _, model := range searchModels {
		indexName := model.IndexName()

		// create or update index
		var index *meilisearch.Index
		if index, err = meilisearchClient.GetIndex(indexName); err != nil {
			if meiliError, ok := err.(meilisearch.Error); ok {
				if meiliError.StatusCode == 404 {
					if _, err = meilisearchClient.CreateIndex(&meilisearch.IndexConfig{
						Uid:        indexName,
						PrimaryKey: model.PrimaryKey(),
					}); err != nil {
						utils.Logger.Panic("Cannot create index "+indexName, zap.Error(err))
					}

					index, err = meilisearchClient.GetIndex(indexName)
					if err != nil {
						utils.Logger.Panic("Cannot get index "+indexName, zap.Error(err))
					}
				}
			} else {
				utils.Logger.Panic("Cannot get index "+indexName, zap.Error(err))
			}
		}

		var filterableAttributes = model.FilterableAttributes()
		if _, err = index.UpdateFilterableAttributes(&filterableAttributes); err != nil {
			utils.Logger.Panic("Cannot update filterable attributes", zap.Error(err))
		}

		var searchableAttributes = model.SearchableAttributes()
		if _, err = index.UpdateSearchableAttributes(&searchableAttributes); err != nil {
			utils.Logger.Panic("Cannot update searchable attributes", zap.Error(err))
		}

		var sortableAttributes = model.SortableAttributes()
		if _, err = index.UpdateSortableAttributes(&sortableAttributes); err != nil {
			utils.Logger.Panic("Cannot update sortable attributes", zap.Error(err))
		}

		var rankingRules = model.RankingRules()
		if _, err = index.UpdateRankingRules(&rankingRules); err != nil {
			utils.Logger.Panic("Cannot update ranking rules", zap.Error(err))
		}
	}
}

type SearchModel interface {
	IDModel
	IndexName() string
	PrimaryKey() string
	FilterableAttributes() []string
	SearchableAttributes() []string
	SortableAttributes() []string
	RankingRules() []string
}

func SearchAddOrReplace[T IDTabler](model T) (err error) {
	if config.Config.MeilisearchUrl == "" {
		return
	}
	_, err = meilisearchClient.Index(model.TableName()).AddDocuments([]T{model}, "id")
	return err
}

func SearchAddOrReplaceInBatch[T IDTabler](models []T) (err error) {
	if config.Config.MeilisearchUrl == "" {
		return
	}
	_, err = meilisearchClient.Index(models[0].TableName()).AddDocuments(models, "id")
	return err
}

func SearchDelete[T IDTabler](id int) (err error) {
	var model T
	if config.Config.MeilisearchUrl == "" {
		return
	}
	_, err = meilisearchClient.Index(model.TableName()).DeleteDocument(strconv.Itoa(id))
	return err
}

func Search[T IDTabler](tx *gorm.DB, models *[]T, q string, filter string, sort []string, columnName string, request utils.PageRequest) (total int, err error) {
	var (
		model     T
		indexName = model.TableName()
		resp      *meilisearch.SearchResponse
	)

	if meilisearchClient == nil {
		return searchFromDB(tx, models, q, filter, sort, columnName, request)
	}

	for i := range sort {
		sort[i] = strings.Replace(sort[i], " ", ":", -1)
	}

	if resp, err = meilisearchClient.Index(indexName).Search(q, &meilisearch.SearchRequest{
		Filter:      filter,
		HitsPerPage: int64(request.PageSize),
		Page:        int64(request.PageNum),
		Sort:        sort,
	}); err != nil {
		return
	}

	total = int(resp.TotalHits)
	if total == 0 {
		return
	}

	// 获取 id 数组
	var idArray []int
	for _, hit := range resp.Hits {
		idArray = append(idArray, int(hit.(Map)["id"].(float64)))
	}

	// 从数据库中读取数据
	err = LoadModelByIDArray(tx, models, idArray)
	return
}

func searchFromDB[T IDTabler](tx *gorm.DB, models *[]T, q string, filter string, sort []string, columnName string, request utils.PageRequest) (total int, err error) {
	// 构造排序字符串
	var orderString strings.Builder
	for i := range sort {
		orderString.WriteString(sort[i])
		if i != len(sort)-1 {
			orderString.WriteString(",")
		}
	}
	if err = request.QuerySet(tx).Where(filter).Order(orderString.String()).
		Where("? like ?", gorm.Expr(columnName), "%"+q+"%").Find(models).Error; err != nil {
		return
	}
	total = len(*models)
	return
}
