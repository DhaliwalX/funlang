package res

import (
    "fmt"
    "strings"
)

type dfsUtil struct {
    dfnum []int
    parent []int
    current int
}

func (u *dfsUtil) init() {
    for i := range u.dfnum {
        u.dfnum[i] = -1
    }
}

func (u *dfsUtil) dfs(g Graph, r int) {
    u.dfnum[r] = u.current

    u.current++
    for _, n := range g[r].Succs {
        if u.dfnum[n] == -1 {
            u.dfs(g, n)
            u.parent[n] = r
        }
    }
}

func dfs(g Graph) *dfsUtil {
    d := dfsUtil{dfnum: make([]int, len(g)), parent: make([]int, len(g)), current: 0}
    d.init()
    d.dfs(g, 0)
    return &d
}

func (u *dfsUtil) String() string {
    builder := strings.Builder{}
    for i, num := range u.dfnum {
        builder.WriteString(fmt.Sprintf("%d - %d\n", i, num))
    }

    return builder.String()
}
