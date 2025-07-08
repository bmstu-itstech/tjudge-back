package service

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/bmstu-itstech/tjudge-back/internal/domain/tjudge"
)

type JudgeLauncher struct {
	judgePath string
}

var ErrNoSheBang = errors.New("shebang missing from file")

func HasSheBang(prog tjudge.Program) error {
	f, err := os.OpenFile(string(prog.Path), os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	var b []byte
	if _, err := f.Read(b); err != nil {
		return err
	}
	if !bytes.HasPrefix(b, []byte("#!")) {
		return ErrNoSheBang
	}
	return nil
}

func (j JudgeLauncher) Run(game tjudge.Game, progs []tjudge.Program) ([]tjudge.RunResult, error) {
	if len(progs) != int(game.Players) {
		return nil, tjudge.ErrInvalidPlayerCount
	}

	for _, p := range progs {
		if err := HasSheBang(p); err != nil {
			return nil, err
		}
	}

	args := []string{string(game.Id)}
	for _, p := range progs {
		args = append(args, string(p.Path))
	}
	judgeCmd := exec.Command(j.judgePath, args...)

	var out bytes.Buffer
	judgeCmd.Stdout = &out

	err := judgeCmd.Run()
	if err != nil {
		// not on unix? too bad!
		// horrible crutches throughout, ancestors cry, this terrible fortune
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil,
				tjudge.NewErrProgramException(fmt.Sprintf("Error in program %d", exitError.ExitCode()))
		}
		return nil, err
	}

	results := make([]tjudge.RunResult, 0)
	for i := range int(game.Players) {
		score_str, err := out.ReadString(' ')
		if err != nil {
			return nil, err
		}
		score_str, _ = strings.CutSuffix(score_str, " ")
		score, err := strconv.Atoi(score_str)
		if err != nil && !(i == int(game.Players)-1 && err == io.EOF) {
			return nil, err
		}
		results = append(results, tjudge.RunResult{
			ProgramId: progs[i].Id,
			Score:     score,
		})
	}
	return results, nil
}
