package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 是一个全局变量，用于存储数据库连接实例
var DB *gorm.DB

// InitMySQL 函数用于初始化 MySQL 连接、创建表以及获取增删改查句柄
func InitMySQL() error {
	// 连接 MySQL 数据库
	dsn := "root:112233@tcp(mysql:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// 将数据库连接实例保存到全局变量中
	DB = db

	// 自动迁移（创建）表结构
	err = AutoMigrate()
	if err != nil {
		return err
	}

	return nil
}

// AutoMigrate 函数用于自动创建表结构
func AutoMigrate() error {
	// 自动迁移（创建）表结构
	err := DB.AutoMigrate(&Artifact{})
	if err != nil {
		return err
	}
	// 可以添加其他表的迁移操作

	return nil
}

// Artifact 是文物信息的模型结构体
type Artifact struct {
	ID                  uint   `gorm:"primary_key;auto_increment"` // 自增主键
	ArtifactID          string `gorm:"unique_index"`               // 文物唯一标识
	Name                string // 文物名称
	Description         string // 文物描述
	Location            string // 文物位置
	Source3DFileURI     string // 文物三维数据的源文件 URI
	Compressed3DFileURI string // 文物三维数据的压缩文件 URI
}

// GetArtifactHandle 函数用于获取文物信息的增删改查句柄
func GetArtifactHandle() *gorm.DB {
	return DB.Model(&Artifact{})
}

// InsertArtifact 函数用于向数据库中插入一条文物信息记录
func InsertArtifact(artifactId, name, description, location, sourceURI, compressedURI string) error {
	artifact := &Artifact{
		ArtifactID:          artifactId,
		Name:                name,
		Description:         description,
		Location:            location,
		Source3DFileURI:     sourceURI,
		Compressed3DFileURI: compressedURI,
	}
	return DB.Model(&Artifact{}).Create(artifact).Error
}

// QueryArtifacts 函数用于查询所有文物信息记录
func QueryArtifacts(artifactId string) ([]Artifact, error) {
	var artifacts []Artifact
	err := DB.Model(&Artifact{}).Where(&Artifact{ArtifactID: artifactId}).Find(&artifacts).Error
	return artifacts, err
}

// UpdateArtifact 函数用于更新文物信息记录
func UpdateArtifact(artifactId, name, description, location string) error {
	updateMap := map[string]interface{}{
		"name":        name,
		"description": description,
		"location":    location,
	}
	err := DB.Model(&Artifact{}).Where(&Artifact{ArtifactID: artifactId}).Updates(updateMap).Error
	return err
}
