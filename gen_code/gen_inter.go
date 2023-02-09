package gen_code

import "modelctl/model"

type IDBMetaData interface {
	ConnectionDB() error
	AllTableData() (model.TableMataDataList, error)
	GetTableColumns(tableName string) (model.ColumnMetaDataList, error)
	TableDataForSelect(flagModel model.FlagModel) (model.TableMataDataList, error)
	//SpecifiedTables(tableName []string) (db_meta_data.TableMetaDataList, error)
}
