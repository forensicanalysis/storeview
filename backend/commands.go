package backend

import "github.com/forensicanalysis/storeview/cobraserver"

func Commands() []*cobraserver.Command {
	return []*cobraserver.Command{ListTables(), SelectItems(), LoadFile(), ListTree(), ListTasks(),
		Files(), Logs(), ErrorsCommand(), Label(), Labels(),
		Query(), RunTask()}
}
