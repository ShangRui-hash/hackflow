package hackflow

type Sqlmap struct {
	name     string
	execPath string
	desp     string
}

func newSqlmap() Tool {
	return &Sqlmap{
		name: SQLMAP,
		desp: "自动化sql注入工具",
	}
}

type SqlmapConfig struct {
	TargetURL   string
	Proxy       string
	BulkFile    string
	RandomAgent bool
	Batch       bool
}

func (s *Sqlmap) Name() string {
	return s.name
}
func (s *Sqlmap) Desp() string {
	return s.desp
}
func (s *Sqlmap) ExecPath() (string, error) {
	return s.execPath, nil
}

func (s *Sqlmap) download() error {
	return nil
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

func NewSqlmap() *Sqlmap {
	tool := container.Get(SQLMAP)
	if tool == nil {
		tool = &Sqlmap{
			name:     SQLMAP,
			execPath: "",
		}
		container.Set(tool)
	}
	return tool.(*Sqlmap)
}
