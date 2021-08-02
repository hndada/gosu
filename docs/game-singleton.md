var (
    g Game
    sSelect SceneSelect
)

init() {
    g = NewGame()
    sSelect = NewSceneSelect()
}


type Game struct {
    cwd string // current working dir
}

func (g Game) GetCWD() string {
    return g.cwd
}
