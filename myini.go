// Package myini --
package myini

import (
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"

	"github.com/alrusov/log"
	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// ServiceCommand --
var ServiceCommand = ""

// IniFileName --
var IniFileName string
var iniFile *ini.File

//----------------------------------------------------------------------------------------------------------------------------//

// OpenIniFile --
func OpenIniFile() bool {
	ret := true
	var err error
	var files []string

	opts := ini.LoadOptions{
		IgnoreInlineComment:        true,
		AllowPythonMultilineValues: true,
	}

	if files, err = filepath.Glob(IniFileName); err != nil {
		log.Message(log.CRIT, "ini: %s", err.Error())
		ret = false
	} else if len(files) == 0 {
		log.Message(log.CRIT, `No ini files "%s" found`, IniFileName)
		ret = false
	} else {
		var filesI []interface{}
		if len(files) > 1 {
			for _, fn := range files[1:] {
				filesI = append(filesI, fn)
			}
		}

		if iniFile, err = ini.LoadSources(opts, files[0], filesI...); err != nil {
			log.Message(log.CRIT, err.Error())
			ret = false
		}
	}

	return ret
}

// GetSectionValues --
func GetSectionValues(name string) misc.StringMap {
	ret := make(misc.StringMap)

	sect := iniFile.Section(name)
	if sect != nil {
		keys := sect.KeyStrings()
		for i := 0; i < len(keys); i++ {
			keyName := keys[i]
			ret[keyName] = iniFile.Section(name).Key(keyName).String()
		}
	}

	return ret
}

// IsParamExists --
func IsParamExists(section string, name string) bool {
	ret := false

	sect := iniFile.Section(section)
	if sect != nil {
		keys := sect.KeyStrings()
		for i := 0; i < len(keys); i++ {
			if keys[i] == name {
				ret = true
				break
			}
		}
	}

	return ret
}

//----------------------------------------------------------------------------------------------------------------------------//

func notFoundError(name string, sectionName ...string) {
	sectList := ""
	for i, sn := range sectionName {
		if i > 0 {
			sectList += ", "
		}
		sectList += `"` + sn + `"`
	}

	log.Message(log.CRIT, `The parameter "%s" is undefined in sections %s`, name, sectList)
}

//----------------------------------------------------------------------------------------------------------------------------//

// GetString --
func GetString(name string, defValue string, mandatory bool, sectionName ...string) string {
	ret := ""
	ok := false

	for _, sn := range sectionName {
		if IsParamExists(sn, name) {
			ret = strings.TrimSpace(iniFile.Section(sn).Key(name).String())
			ok = true
			break
		}
	}

	if !ok {
		if mandatory {
			notFoundError(name, sectionName...)
		} else {
			ret = defValue
		}
	}

	return ret
}

// GetInt --
func GetInt(name string, defValue int, mandatory bool, sectionName ...string) int {
	ret := 0
	badValue := -999999999
	ok := false

	for _, sn := range sectionName {
		if IsParamExists(sn, name) {
			ret = iniFile.Section(sn).Key(name).MustInt(badValue)
			ok = true
			break
		}
	}

	if !ok {
		if mandatory {
			notFoundError(name, sectionName...)
		} else {
			ret = defValue
		}
	}

	if ret == badValue {
		log.Message(log.CRIT, "The parameter \"%s\"  has bad value (must be int)", name)
	}

	return ret
}

// GetBool --
func GetBool(name string, defValue bool, mandatory bool, sectionName ...string) bool {
	ret := false
	ok := false

	for _, sn := range sectionName {
		if IsParamExists(sn, name) {
			ret = iniFile.Section(sn).Key(name).MustBool(defValue)
			ok = true
			break
		}
	}

	if !ok {
		if mandatory {
			notFoundError(name, sectionName...)
		} else {
			ret = defValue
		}
	}

	return ret
}

// GetFloat --
func GetFloat(name string, defValue float64, mandatory bool, sectionName ...string) float64 {
	ret := float64(0)
	badValue := float64(-999999999999)
	ok := false

	for _, sn := range sectionName {
		if IsParamExists(sn, name) {
			ret = iniFile.Section(sn).Key(name).MustFloat64(badValue)
			ok = true
			break
		}
	}

	if !ok {
		if mandatory {
			notFoundError(name, sectionName...)
		} else {
			ret = defValue
		}
	}

	if ret == badValue {
		log.Message(log.CRIT, "The parameter \"%s\"  has bad value (must be float)", name)
	}

	return ret
}

//----------------------------------------------------------------------------------------------------------------------------//
