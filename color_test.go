package oiio

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	OCIO_CONFIG_PATH string
)

func init() {
	// Set the environment to file-based test data config
	pwd, _ := os.Getwd()
	pwd, _ = filepath.Abs(pwd)

	OCIO_CONFIG_PATH = filepath.Join(pwd, "testdata/spi-vfx/config.ocio")
	if _, err := os.Stat(OCIO_CONFIG_PATH); os.IsNotExist(err) {
		panic(fmt.Sprintf("%s not found\n", OCIO_CONFIG_PATH))
	} else {
		os.Setenv("OCIO", OCIO_CONFIG_PATH)
	}
}

func TestNewColorConfig(t *testing.T) {
	if !SupportsOpenColorIO() {
		t.SkipNow()
	}
	cfg, err := NewColorConfig()
	if err != nil {
		t.Fatal(err)
	}
	checkConfig(t, cfg)
}

func TestNewColorConfigPath(t *testing.T) {
	if !SupportsOpenColorIO() {
		t.SkipNow()
	}
	_, err := NewColorConfigPath(OCIO_CONFIG_PATH)
	if err != nil {
		t.Error(err)
	}
}

func checkConfig(t *testing.T, cfg *ColorConfig) {
	expected := 20
	actual := cfg.NumColorSpaces()
	if expected != actual {
		t.Errorf("Expected number of colorspaces == %v; got %v", expected, actual)
	}

	expected = 1
	actual = cfg.NumLooks()
	if expected != actual {
		t.Errorf("Expected number of looks == %v; got %v", expected, actual)
	}

	expected = 2
	actual = cfg.NumDisplays()
	if expected != actual {
		t.Errorf("Expected number of displays == %v; got %v", expected, actual)
	}

	expected = 3
	actual = cfg.NumViews("sRGB")
	if expected != actual {
		t.Errorf("Expected number of views for sRGB == %v; got %v", expected, actual)
	}

	cs_expect := "vd8"
	cs_actual := cfg.ColorSpaceNameByRole("matte_paint")
	if cs_expect != cs_actual {
		t.Errorf("Expected ColorSpace for role 'matte_paint' == %v; got %v", cs_expect, cs_actual)
	}

	for i, cs_expect := range []string{"lnf", "lnh", "ln16", "lg16"} {
		cs_actual = cfg.ColorSpaceNameByIndex(i)
		if cs_expect != cs_actual {
			t.Errorf("Expected ColorSpace for index %d == %v; got %v", i, cs_expect, cs_actual)
		}
	}

	look_expect := "di"
	look_actual := cfg.LookNameByIndex(0)
	if look_actual != look_expect {
		t.Errorf("Expected Look for index 0 == %v; got %v", look_expect, look_actual)
	}

	for i, disp_expect := range []string{"DCIP3", "sRGB"} {
		disp_actual := cfg.DisplayNameByIndex(i)
		if disp_expect != disp_actual {
			t.Errorf("Expected Display for index %d == %v; got %v", i, disp_expect, disp_actual)
		}
	}

	for i, view_expect := range []string{"Film", "Log", "Raw"} {
		view_actual := cfg.ViewNameByIndex("sRGB", i)
		if view_actual != view_expect {
			t.Errorf("Expected Look for display sRGB at index %d == %v; got %v", i, view_expect, view_actual)
		}
	}

}

func TestNewColorProcessor(t *testing.T) {
	if !SupportsOpenColorIO() {
		t.SkipNow()
	}
	cfg, err := NewColorConfig()
	if err != nil {
		t.Fatal(err)
	}

	_, err = cfg.CreateColorProcessor("lnf", "vd8")
	if err != nil {
		t.Error(err.Error())
	}
}
