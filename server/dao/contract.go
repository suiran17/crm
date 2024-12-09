package dao

import (
	"time"

	"crm/global"
	"crm/models"
)

type ContractDao struct {
}

func NewContractDao() *ContractDao {
	return &ContractDao{}
}

func (c *ContractDao) Create(param *models.ContractCreateParam) error {
	contract := models.Contract{
		Name:        param.Name,
		Amount:      param.Amount,
		BeginTime:   param.BeginTime,
		OverTime:    param.OverTime,
		Remarks:     param.Remarks,
		Cid:         param.Cid,
		Productlist: param.Productlist,
		Status:      param.Status,
		Creator:     param.Creator,
		Created:     time.Now().Unix(),
	}
	return global.Db.Create(&contract).Error
}

func (c *ContractDao) Update(param *models.ContractUpdateParam) error {
	contract := models.Contract{
		Id:          param.Id,
		Name:        param.Name,
		Amount:      param.Amount,
		BeginTime:   param.BeginTime,
		OverTime:    param.OverTime,
		Remarks:     param.Remarks,
		Cid:         param.Cid,
		Productlist: param.Productlist,
		Status:      param.Status,
		Updated:     time.Now().Unix(),
	}
	db := global.Db.Model(&contract).Select("*").Omit("id", "creator", "created")
	return db.Updates(&contract).Error
}

func (c *ContractDao) Delete(param *models.ContractDeleteParam) error {
	return global.Db.Delete(&models.Contract{}, param.Ids).Error
}

// GetList 根据查询参数获取合同列表，并返回总记录数和可能的错误。
// param: 合同查询参数，包含分页信息、过滤条件等。
// 返回值: 合同列表数组，总记录数，以及可能的错误。
func (c *ContractDao) GetList(param *models.ContractQueryParam) ([]*models.ContractList, int64, error) {
	// 初始化合同列表数组
	contractList := make([]*models.ContractList, 0)
	// 定义查询的字段
	field := "contract.id, contract.name, contract.amount, contract.begin_time, contract.over_time, customer.name as cname, contract.remarks, contract.status, contract.created, contract.updated"
	// 定义JOIN条件，基于合同创建者和客户表关联
	where := "inner join customer on contract.cid = customer.id and contract.creator = ?"
	// 定义计数查询的原始SQL
	raw := "select count(*) from contract where creator = ?"

	// 计算分页查询的偏移量
	// 分页查询
	offset := (param.Page.PageNum - 1) * param.Page.PageSize
	// 构建基础查询，包括分页和选择字段
	db := global.Db.Offset(offset).Limit(param.Page.PageSize).Table(CONTRACT).Select(field)

	var rows int64
	// 根据查询参数是否存在ID或状态，动态构建查询条件并计数
	if param.Id != NumberNull {
		// 如果提供了ID，则在JOIN条件后添加ID过滤条件，并重新计算记录数
		db.Joins(where+" and contract.id = ?", param.Creator, param.Id)
		global.Db.Raw(raw+" and contract.id = ?", param.Creator, param.Creator).Scan(&rows)
	} else if param.Status != NumberNull {
		// 如果未提供ID但提供了状态，则在JOIN条件后添加状态过滤条件，并重新计算记录数
		db.Joins(where+" and contract.status = ?", param.Creator, param.Status)
		global.Db.Raw(raw+" and contract.status = ?", param.Creator, param.Status).Scan(&rows)
	} else {
		// 如果既未提供ID也未提供状态，仅应用JOIN条件，并计算记录数
		db.Joins(where, param.Creator)
		global.Db.Raw(raw, param.Creator).Scan(&rows)
	}
	// 执行查询并填充合同列表数组，同时检查是否有错误发生
	if err := db.Scan(&contractList).Error; err != nil {
		// 如果查询出错，返回空列表、特殊的NULL值和错误
		return nil, NumberNull, nil
	}
	// 返回查询结果列表、总记录数和无错误
	return contractList, rows, nil
}

func (c *ContractDao) GetListByUid(uid int64) ([]*models.ContractList, error) {
	contracts := make([]*models.ContractList, 0)
	s := "contract.id, contract.name, contract.amount, contract.begin_time, contract.over_time, customer.name as cname, contract.remarks, contract.status, contract.created, contract.updated"
	j := "left join customer on contract.cid = customer.id and contract.creator = ?"
	err := global.Db.Table(CONTRACT).Select(s).Joins(j, uid).Scan(&contracts).Error
	if err != nil {
		return nil, err
	}
	return contracts, nil
}

func (c *ContractDao) GetInfo(param *models.ContractQueryParam) (*models.ContractInfo, error) {
	contract := models.Contract{
		Id: param.Id,
	}
	contractInfo := models.ContractInfo{}
	err := global.Db.Table(CONTRACT).Where(&contract).First(&contractInfo).Error
	if err != nil {
		return nil, err
	}
	return &contractInfo, nil
}

func (c *ContractDao) GetAddedPList(id int64) (*models.Contract, error) {
	var contract models.Contract
	err := global.Db.Table(CONTRACT).Select("productlist").First(&contract, id).Error
	if err != nil {
		return nil, err
	}
	return &contract, nil
}
