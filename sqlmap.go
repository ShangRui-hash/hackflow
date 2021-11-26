package hackflow

type Sqlmap struct {
	BaseTool
}

func newSqlmap() Tool {
	return &Sqlmap{
		BaseTool: BaseTool{
			name: SQLMAP,
			desp: "自动化sql注入工具",
		},
	}
}

func GetSqlmap() *Sqlmap {
	return container.Get(SQLMAP).(*Sqlmap)
}

type SqlmapConfig struct {
	TargetURL   string
	Proxy       string
	BulkFile    string
	RandomAgent bool
	Batch       bool
}

func (s *Sqlmap) ExecPath() (string, error) {
	return s.BaseTool.ExecPath(s.Download)
}

func (s *Sqlmap) Download() (string, error) {
	return "", nil
}

//SetDebug 是否开启Debug
func (u *Sqlmap) SetDebug(isDebug bool) {

}
func (s *Sqlmap) Run(config *SqlmapConfig) (resultCh chan string, err error) {
	//todo
	return nil, nil
}

func (s *Sqlmap) GetDbs(config *SqlmapConfig) (dbs []string, err error) {
	//todo
	return nil, nil
}

func (s *Sqlmap) GetTables() (tables []string, err error) {
	//todo
	return nil, nil
}
