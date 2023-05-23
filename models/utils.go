package models

import (
	"ChatDanBackend/utils"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type IDModel interface {
	GetID() int
}

type Tabler interface {
	TableName() string
}

type IDTabler interface {
	IDModel
	Tabler
}

func CacheName[T IDTabler](model T) (name string) {
	return model.TableName() + ":" + strconv.Itoa(model.GetID())
}

func CacheNameFromTableName(tableName string, id int) (name string) {
	return tableName + ":" + strconv.Itoa(id)
}

// PageLoad 分页查询
// tx 数据库包含表和查询条件
func PageLoad[T IDTabler](tx *gorm.DB, models *[]T, key string, request utils.PageRequest) (version, total int, err error) {
	var idArray []int

	// 设置版本号
	if request.Version != 0 {
		version = request.Version
	} else {
		// 读取最新版本号的缓存
		if version, idArray, err = loadLatestVersion(tx, key); err != nil {
			return
		}
		goto PROCESS
	}

	// 从版本号缓存中读取 idArray
	if err = utils.Get(key+":"+strconv.Itoa(version), &idArray); err != nil {
		if err != utils.ErrCacheMiss {
			return
		}

		// 缓存中没有，加载最新版本号的缓存
		if version, idArray, err = loadLatestVersion(tx, key); err != nil {
			return
		}
	}

PROCESS:
	singleTx := tx.Session(&gorm.Session{NewDB: true})

	// 设置总数
	total = len(idArray)
	if total == 0 {
		return
	}

	// 设置分页
	size := request.PageSize
	offset := (request.PageNum - 1) * request.PageSize
	if offset >= total {
		return
	}
	if offset+size > total {
		size = total - offset
	}

	err = LoadModelByIDArray(singleTx, models, idArray[offset:offset+size])

	return
}

type PageVersionValue struct {
	Version int   `json:"version"`
	IDArray []int `json:"id_array"`
}

func loadLatestVersion(tx *gorm.DB, key string) (version int, idArray []int, err error) {
	// 读取最新版本号的缓存
	var value PageVersionValue
	if err = utils.Get(key+":latest", &value); err != nil {
		if err != utils.ErrCacheMiss {
			return
		}

		return SetLatestVersion(tx, key)
	} else {
		idArray = value.IDArray
		version = value.Version
	}
	return
}

func SetLatestVersion(tx *gorm.DB, key string) (version int, idArray []int, err error) {
	// 生成当前版本号
	version = int(time.Now().UnixMicro())

	// 缓存中没有，从数据库中加载 idArray
	if err = tx.Pluck("id", &idArray).Error; err != nil {
		return
	}
	if len(idArray) == 0 {
		return
	}

	// 设置版本号缓存
	if err = utils.Set(key+":"+strconv.Itoa(version), idArray, 10*time.Minute); err != nil {
		return
	}

	// 设置最新版本号缓存
	// 这里设置的过期时间要比版本号缓存的过期时间要短，防止其他请求读到错误的最新版本号
	if err = utils.Set(key+":latest", PageVersionValue{
		Version: version,
		IDArray: idArray,
	}, 9*time.Minute); err != nil {
		return
	}

	return
}

// LoadModel 从数据库或缓存加载数据
func LoadModel[T IDTabler](tx *gorm.DB, model *T) (err error) {
	// 先从缓存中加载
	if err = utils.Get(CacheName(*model), model); err != nil {
		if err != nil {
			if err != utils.ErrCacheMiss {
				return
			}

			// 缓存中没有，从数据库中加载
			if err = tx.First(model).Error; err != nil {
				return
			}

			// 设置缓存
			if err = utils.Set(CacheName(*model), model, 10*time.Minute); err != nil {
				return
			}
		}
	}
	return err
}

func LoadModelByIDArray[T IDTabler](tx *gorm.DB, models *[]T, idArray []int) (err error) {
	var (
		_model    T
		tableName = _model.TableName()
	)
	// 设置总数
	size := len(idArray)
	if size == 0 {
		return
	}

	// 从缓存或数据库中加载 models
	// 构建查询数据
	*models = make([]T, size)
	notCachedModels := make([]T, 0, size)        // 未缓存的数据
	notCachedIdArray := make([]int, 0, size)     // 未缓存的数据的 id
	notCacheIdMapping := make(map[int]int, size) // 未缓存的数据的 id 与 models 的索引映射

	// 从缓存中读取数据
	for i, id := range idArray {
		name := CacheNameFromTableName(tableName, id)
		if err = utils.Get(name, &(*models)[i]); err != nil {
			if err == utils.ErrCacheMiss {
				notCachedIdArray = append(notCachedIdArray, id)
				notCacheIdMapping[id] = i
				continue
			}
			return
		}
	}

	// 从数据库中读取没有缓存的数据
	if len(notCachedIdArray) > 0 {
		if err = tx.Where("id IN ?", notCachedIdArray).Find(&notCachedModels).Error; err != nil {
			return
		}

		// 将数据放入缓存
		for i := range notCachedModels {
			name := CacheName(notCachedModels[i])
			if err = utils.Set(name, notCachedModels[i], 10*time.Minute); err != nil {
				return
			}
		}

		// 根据 id 与 models 的索引映射，将数据放入 models 中
		for i := range notCachedModels {
			(*models)[notCacheIdMapping[notCachedModels[i].GetID()]] = notCachedModels[i]
		}

		// 检查是否所有数据都读取到了
		loadedModels := make([]T, 0, size)
		for _, model := range *models {
			if model.GetID() != 0 {
				loadedModels = append(loadedModels, model)
			}
		}

		if len(loadedModels) != size {
			*models = loadedModels
		}
	}

	return
}

func LoadModelAll[T IDTabler](tx *gorm.DB, model *[]T) (err error) {
	var (
		_model    T
		tableName = _model.TableName()
	)

	// 从缓存中读取数据
	if err = utils.Get(tableName, model); err != nil {
		if err != utils.ErrCacheMiss {
			return
		}

		// 从数据库中读取数据
		if err = tx.Find(model).Error; err != nil {
			return
		}

		// 将数据放入缓存
		if err = utils.Set(tableName, model, 10*time.Minute); err != nil {
			return
		}
	}
	return
}

func CreateModel[T IDTabler](tx *gorm.DB, model *T) (err error) {
	if err = tx.FirstOrCreate(model).Error; err != nil {
		return
	}

	// 设置缓存
	if err = utils.Set(CacheName(*model), model, 10*time.Minute); err != nil {
		return
	}

	return
}

func UpdateModel[T IDTabler](tx *gorm.DB, model *T, columns any) (err error) {
	if err = tx.Model(model).Updates(columns).Error; err != nil {
		return
	}

	// 设置缓存
	if err = utils.Set(CacheName(*model), model, 10*time.Minute); err != nil {
		return
	}

	return
}

func DeleteModel[T IDTabler](tx *gorm.DB, model *T) (err error) {
	if err = tx.Delete(model).Error; err != nil {
		return
	}

	// 删除缓存
	utils.Delete(CacheName(*model))

	return
}
