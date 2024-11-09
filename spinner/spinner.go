package spinner

import (
	"github.com/charmbracelet/huh/spinner"
)

func SpinnerClass(desc string) {
	_ = spinner.New().Title(desc).Run()
}
