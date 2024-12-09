package dao

import (
	"context"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"crm/common"
	"crm/global"
	"crm/models"
)

const (

	// 数据库表名
	USER        = "user"
	CUSTOMER    = "customer"
	CONTRACT    = "contract"
	PRODUCT     = "product"
	SUBSCRIBE   = "subscribe"
	NOTICE      = "notice"
	MAIL_CONFIG = "mail_config"

	// 空值
	NumberNull = 0
	StringNull = ""

	// 运行环境
	Dev  = "dev"
	Prod = "prod"
)

var ctx = context.Background()

// RestPage 分页查询
// page  设置起始页、每页条数,
// name  查询目标表的名称
// query 查询条件,
// dest  查询结果绑定的结构体,
// bind  绑定表结构对应的结构体
func restPage(page models.Page, name string, query interface{}, dest interface{}, bind interface{}) (int64, error) {
	if page.PageNum > 0 && page.PageSize > 0 {
		offset := (page.PageNum - 1) * page.PageSize
		global.Db.Offset(offset).Limit(page.PageSize).Table(name).Where(query).Find(dest)
	}
	res := global.Db.Table(name).Where(query).Find(bind)
	return res.RowsAffected, res.Error
}

type CommonDao struct {
}

func NewCommonDao() *CommonDao {
	return &CommonDao{}
}

func (c *CommonDao) InitDatabase() error {
	env := global.Config.Server.Runenv
	if env == Prod {
		dbFile := global.Config.Mysql.DbFile
		sql, err := os.ReadFile(dbFile)
		if err != nil {
			log.Printf("Common.InitDatabase.Error: read file %s error: %s", dbFile, err)
			return err
		}
		sqls := strings.Split(string(sql), ";")
		for _, sql := range sqls {
			s := strings.TrimSpace(sql)
			if s == StringNull {
				continue
			}
			if err := global.Db.Exec(s).Error; err != nil {
				log.Printf("Common.InitDatabase.Error: %s", err)
				return err
			}
		}
		return nil
	}
	return nil
}

// FileUpload 实现文件上传功能。
// 它接收一个multipart.FileHeader指针作为参数，返回一个models.FileInfo指针和一个错误。
// 该方法主要负责生成文件的唯一标识，计算文件存储路径，打开并复制文件内容到目标路径，最后返回文件信息。
func (c *CommonDao) FileUpload(file *multipart.FileHeader) (*models.FileInfo, error) {
	// 从全局配置中获取文件存储路径
	dist := global.Config.File.Path
	// 生成唯一文件名，加上原文件名的扩展名
	name := common.GenUUID() + path.Ext(file.Filename)
	// 拼接得到文件的存储路径
	dn := dist + name

	// 打开上传的文件，准备进行复制
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	// 确保文件关闭，避免资源泄露
	defer src.Close()

	// 创建目标文件，准备写入
	out, err := os.Create(dn)
	if err != nil {
		return nil, err
	}
	// 确保文件关闭，避免资源泄露
	defer out.Close()

	// 将上传的文件内容复制到目标文件
	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	// 构建文件信息对象，包含文件的URL和名称
	flieInfo := models.FileInfo{
		Url:  dn,
		Name: name,
	}
	return &flieInfo, err
}

func (c *CommonDao) FileRemove(fileName string) error {
	file := global.Config.File.Path + fileName
	return os.Remove(file)
}
