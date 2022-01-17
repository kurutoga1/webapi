package config_test

import (
	"log"
	"path/filepath"
	"strings"
	"testing"
	"webapi/server/config"
	"webapi/utils"
)

func init() {
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	config.SetProConfPath("programConfig_test.json")
}

func TestProgramConfig_ToProperPath(t *testing.T) {
	p := config.NewProgramConfig()
	p.ProHelpPath = "programs/xxxxx/help.txt"
	p.ProCommand = "python3 programs/xxxxx/eg.py"
	p.ToProperPath()

	if !filepath.IsAbs(p.ProHelpPath) {
		t.Errorf("p.ProHelpPath is not abs.")
	}

	if !strings.Contains(p.ProCommand, currentDir) {
		t.Errorf("p.ProCommand doesn't contain abspath.")
	}
}

func TestProgramConfig_ReplacedCmd(t *testing.T) {
	p := config.NewProgramConfig()
	p.ProCommand = "INPUTFILE OUTPUTDIR PARAMETA"
	cmd := p.ReplacedCmd("i", "o", "p")

	if cmd != "i o p" {
		t.Errorf("ReplaceCmd(): %v, want: i o p \n", cmd)
	}
}

func TestGetProConfByName(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		testGetProConfByName(t, "convertToJson", true)
	})
	t.Run("test 2", func(t *testing.T) {
		testGetProConfByName(t, "dummy", false)
	})
}

func testGetProConfByName(t *testing.T, programName string, want bool) {
	_, err := config.GetProConfByName(programName)
	if (err != nil) == want {
		t.Errorf(err.Error())
	}
}

func TestGetPrograms(t *testing.T) {
	programConfigHolders, err := config.GetPrograms()
	if err != nil {
		t.Errorf("err from GetPrograms(): %v \n", err.Error())
	}
	for _, p := range programConfigHolders {
		if p.Name() != "convertToJson" {
			t.Errorf("p.Name() is not %v \n", "convertToJson")
		}
	}
}
