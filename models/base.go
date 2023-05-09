package models

import "github.com/jinzhu/copier"

type Map = map[string]any

var CopyOption = copier.Option{IgnoreEmpty: true}
