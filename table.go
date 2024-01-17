package main

import (
	"fmt"
	"html/template"
	"io"
	"math"
)

func Print(results []Result, names map[uint]string, w io.Writer) error {
	t, err := template.New("Template").Funcs(template.FuncMap{
		"Name": func(sailNumber uint) string {
			return names[sailNumber]
		},
		"Inc": func(x int) int {
			return x + 1
		},
		"IsGrey": func(line int) bool {
			return line%2 == 1
		},
		"ScoreOf": func(score float32) string {
			if math.Mod(float64(score), 1.0) < 0.01 {
				return fmt.Sprintf("%.0f", score)
			} else {
				return fmt.Sprintf("%.1f", score)
			}
		},
	}).Parse(templ)
	if err != nil {
		return err
	}

	return t.Execute(w, results)
}

var templ = `<table style="margin-left:auto;margin-right:auto">
  <tr>
    <td style="padding:5px"><b>Sail</b></td>
    <td style="padding:5px"><b>Sailor</b></td>
    <td style="padding:5px"><b>Club</b></td>
    <td style="padding:5px"><b>Design</b></td>
    <td style="padding:5px"><b>Rank</b></td>
    <td style="padding:5px"><b>Tot</b></td>
{{- range $i, $race := (index . 0).Scores}}
    <td style="padding:5px;text-align:center"><b>R {{Inc $i}}</b></td>
{{- end -}}
  </tr>
{{- range $i, $result := .}}
  <tr{{if IsGrey $i}} style="background-color: #dddddd"{{end}}>
    <td style="padding:5px"><b>{{$result.SailNumber}}</b></td>
    <td style="padding:5px"><b>{{Name $result.SailNumber}}</b></td>
    <td style="padding:5px"><b>SARYC</b></td>
    <td style="padding:5px"><b></b></td>
    <td style="padding:5px;text-align:center"><b>{{$result.Rank}}</b></td>
    <td style="padding:5px;text-align:center"><b>{{ScoreOf $result.Tot}}</b></td>
  {{- range $j, $score := $result.Scores}}
    <td style="text-align:center{{if $score.Drop}};background-color:grey;color:white{{end}}"><b>{{ScoreOf $score.Score}}</b></td>
  {{- end}}
  </tr>
{{- end}}
</table>`
