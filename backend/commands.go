package backend

import "github.com/forensicanalysis/storeview/cobraserver"

var Commands = []*cobraserver.Command{ListTables(), SelectItems(), LoadFile(), ListTree(), ListTasks(),
	Files(), Logs(), ErrorsCommand(), Label(), Labels(),
	Query(), RunTask()}
