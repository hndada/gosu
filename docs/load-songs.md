디렉토리 구조
C:\Music1
ㄴ`<osu-song1>`
ㄴ`<o2jam-song1>`
ㄴ`<osu-song2>`
ㄴ`<osu-song3>`
ㄴ`<o2jam-song2>`
D:\Music2
ㄴ`<osu-song4>`
ㄴ`<o2jam-song3>`

// default:  수정 시간이 갱신 시간 이후 (aka 이후 변경된 것)만 갱신
// forced: 전부 갱신
// 1. Music을 통째로 스킵하려면 마지막 갱신 시간을 별도로 저장해둬야 함
    // 임시로 게임 실행할 때마다 forced load(-> 초기화 때 마지막 수정 시간을 0으로 해두면 됨); singleton에 저장 
// 2. Music을 일단 한번 스캔하려 한다면 Music 폴더의 수정 시간을 사용

// Select Scene으로 넘어올 때 언제나 실행
func (s *SceneSelect) LoadCharts(forced bool) {
    if !forced && Music 폴더가 최신일 때 {
        goto mark
    }
    for i, dir := range dirs { // chart directories
        if !forced && ChartSet 폴더가 최신일 때 {
            continue
        }
        var cs []ChartInfo
        for j, f := range dir { // chart in a chart directory
            var c *Chart
            var err error
            filepath:=path.Abs(f)
            switch '확장자명'{
            case '.osu':
                switch game.mode(filepath) { // return value: Mode 코드값
                case game.Mania:
                    if c, err = mania.NewChartFromOsu(filepath); err != nil {
                        panic(err)
                        // continue
                    }
                default:
                    continue
                }
            case '.ojn':
                continue
            }
            append(cs, ChartInfo{
                Path: filepath
                Header: c.Header 
                Level: c.Level // 모드 구현 시 모드 별 level 미리 계산
            })
        }
        append(s.ChartSets, cs)
    }
    sort ChartSets by 이름/난이도
    mark:
        갱신 시간 업데이트
}

// 내부에서 osu.Parse(filepath) 호출
// 얘는 osu.Format에 대해서 알 필요가 X
func NewChartFromOsu(path string) (*Chart, error) { // package mania
	var c Chart
    o := osu.Parse(path)
    c.Header = game.ParseHeaderFromOsu(o)
    c.Keys = int(c.Parameter["Scale"])
	if err := c.loadNotes(o); err != nil {
		return nil, err
	}
	c.CalcDifficulty()
	return &c, nil
}

type ChartInfo struct {
    MusicDirNo int // Music directory number
    Path string
    Header ChartHeader
    Level float64 // 모드 구현 시 모드 별 level 미리 계산
    (SetName string) // Header에 있을텐데, 무결점 보장이 안돼있음
}