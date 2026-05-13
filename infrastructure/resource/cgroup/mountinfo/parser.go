package mountinfo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Compogo/resource/infrastructure/resource/cgroup/path"
	"github.com/Compogo/types/set"
	"github.com/spf13/cast"
)

const (
	Separator              = " "
	OptionsSeparator       = ","
	OptionalFieldSeparator = "-"

	IdFieldNumber       uint8 = 0
	ParentIdFieldNumber uint8 = 1
	DeviceIdFieldNumber uint8 = 2
	RootFieldNumber     uint8 = 3
	PointFieldNumber    uint8 = 4
	OptionsFieldNumber  uint8 = 5
	OptionalFieldNumber uint8 = 6

	FSTypeFieldNumber        uint8 = 0
	SourceFieldNumber        uint8 = 1
	SupperOptionsFieldNumber uint8 = 2

	CountFirstFields   uint8 = 7
	CountSecondFields  uint8 = 3
	MinimumCountFields uint8 = CountFirstFields + CountSecondFields
)

type Parser struct {
	path   path.Path
	points set.Set[*Point]
}

func NewParser(options ...Option) *Parser {
	p := &Parser{}

	options = append([]Option{WithPath(path.CGroupMountInfoSelfPath)}, options...)

	for _, option := range options {
		option(p)
	}

	return p
}

func (parser *Parser) Points() (set.Set[*Point], error) {
	if parser.points.Len() == 0 {
		if err := parser.Parse(); err != nil {
			return nil, err
		}
	}

	return parser.points, nil
}

func (parser *Parser) Parse() error {
	parser.points.Reset()

	mountinfoFile, err := os.OpenFile(string(parser.path), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer mountinfoFile.Close()

	lineScanner := bufio.NewScanner(mountinfoFile)
	lineScanner.Split(bufio.ScanLines)

	for lineScanner.Scan() {
		point, err := parser.pointFromLine(lineScanner.Text())
		if err != nil {
			return err
		}

		parser.points.Add(point)
	}

	return lineScanner.Err()
}

func (parser *Parser) pointFromLine(line string) (point *Point, err error) {
	fields := strings.Split(line, Separator)

	if len(fields) < int(MinimumCountFields) {
		return nil, fmt.Errorf(
			"mountinfo: new from string '%s' failed, wrong number of fields, minimum expected %d, found %d",
			line,
			MinimumCountFields,
			len(fields),
		)
	}

	id, err := cast.ToUint32E(fields[IdFieldNumber])
	if err != nil {
		return nil, fmt.Errorf("mountinfo: id parse failed: %w", err)
	}

	parentId, err := cast.ToUint32E(fields[ParentIdFieldNumber])
	if err != nil {
		return nil, fmt.Errorf("mountinfo: parent id parse failed: %w", err)
	}

	for index, fieldValue := range fields[OptionalFieldNumber:] {
		if fieldValue == OptionalFieldSeparator {
			fsTypeStartIndex := OptionalFieldNumber + uint8(index) + 1

			if len(fields) != int(fsTypeStartIndex+CountSecondFields) {
				return nil, fmt.Errorf(
					"mountinfo: new from string '%s' failed, wrong number of fields, minimum expected %d, found %d",
					line,
					fsTypeStartIndex+CountSecondFields,
					len(fields),
				)
			}

			mInfoFSTypeFieldNumber := FSTypeFieldNumber + fsTypeStartIndex
			mInfoMountSourceFieldNumber := SourceFieldNumber + fsTypeStartIndex
			mInfoSupperOptionsFieldNumber := SupperOptionsFieldNumber + fsTypeStartIndex

			return &Point{
				Id:             id,
				ParentID:       parentId,
				DeviceID:       fields[DeviceIdFieldNumber],
				Root:           path.Path(fields[RootFieldNumber]),
				MountPoint:     path.Path(fields[PointFieldNumber]),
				Options:        set.NewSet(strings.Split(fields[OptionsFieldNumber], OptionsSeparator)...),
				OptionalFields: set.NewSet(fields[OptionalFieldNumber:(fsTypeStartIndex - 1)]...),
				FSType:         fields[mInfoFSTypeFieldNumber],
				MountSource:    fields[mInfoMountSourceFieldNumber],
				SuperOptions:   set.NewSet(strings.Split(fields[mInfoSupperOptionsFieldNumber], OptionsSeparator)...),
			}, nil
		}
	}

	return nil, fmt.Errorf("mountinfo: line '%s' parse failed, invalid format", line)
}
