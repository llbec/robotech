package dfusesquadron_test

import (
	"testing"

	"github.com/llbec/robotech/squadrons/dfusesquadron"
)

func Test_ReadBlock(t *testing.T) {
	st := dfusesquadron.NewStream("server_7d59808b3e874fd899e7f785e546642a", "heco.streamingfast.io:443", "true")
	st.Search(-10, 0)
}
