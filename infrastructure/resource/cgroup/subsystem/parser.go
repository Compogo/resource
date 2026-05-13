package subsystem

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Compogo/resource/infrastructure/resource/cgroup/path"
	"github.com/Compogo/types/linker"
	"github.com/spf13/cast"
)

const (
	separator          = ":"
	subsystemSeparator = ","

	IdFieldNumber            uint8 = 0
	subsystemFieldNumber     uint8 = 1
	mountInfoPathFieldNumber uint8 = 2

	countFields uint8 = 3
)

type Parser struct {
	path       path.Path
	subsystems *linker.Linker[Name, *Subsystem]
}

func NewParser(options ...Option) *Parser {
	p := &Parser{
		subsystems: linker.NewLinker[Name, *Subsystem](),
	}

	options = append([]Option{WithPath(path.CGroupSubsystemSelfPath)}, options...)

	for _, option := range options {
		option(p)
	}

	return p
}

func (parser *Parser) Subsystems() (*linker.Linker[Name, *Subsystem], error) {
	if parser.subsystems.Len() == 0 {
		if err := parser.Parse(); err != nil {
			return nil, err
		}
	}

	return parser.subsystems, nil
}

func (parser *Parser) Parse() error {
	parser.subsystems.Reset()

	cgroupSubsystemFile, err := os.OpenFile(string(parser.path), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer cgroupSubsystemFile.Close()

	lineScanner := bufio.NewScanner(cgroupSubsystemFile)
	lineScanner.Split(bufio.ScanLines)

	var subsystems []*Subsystem
	for lineScanner.Scan() {
		subsystems, err = parser.subsystemFromLine(lineScanner.Text())
		if err != nil {
			return err
		}

		for _, subsystem := range subsystems {
			parser.subsystems.Add(subsystem.Name, subsystem)
		}
	}

	return lineScanner.Err()
}

func (parser *Parser) subsystemFromLine(line string) ([]*Subsystem, error) {
	fields := strings.SplitN(line, separator, int(countFields))

	if len(fields) != int(countFields) {
		return nil, fmt.Errorf(
			"cgroup.subsystem: new from string '%s' failed, wrong number of fields, expected %d, found %d",
			line,
			countFields,
			len(fields),
		)
	}

	id, err := cast.ToUint8E(fields[IdFieldNumber])
	if err != nil {
		return nil, fmt.Errorf("cgroup.subsystem: subsystem id parse failed: %w", err)
	}

	subsystemNames := strings.Split(fields[subsystemFieldNumber], subsystemSeparator)
	subsystems := make([]*Subsystem, len(subsystemNames))

	for index, subSystemName := range subsystemNames {
		name := Name(subSystemName)
		if !AllNames.Contains(name) {
			return nil, fmt.Errorf("cgroup.subsystem: subsystem '%s' has invalid name", name)
		}

		subsystems[index] = &Subsystem{
			Id:                 id,
			Name:               name,
			MountPointRootPath: path.Path(fields[mountInfoPathFieldNumber]),
		}
	}

	return subsystems, nil
}
