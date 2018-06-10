package sheetsmap

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/Iwark/spreadsheet"
	"github.com/grokify/gotilla/strings/stringsutil"
)

const (
	ErrorColumnNotFound = "ErrorColumnNotFound"
	ErrorEnumNotMatched = "ErrorEnumNotMatched"
)

type SheetsMap struct {
	GoogleClient   *http.Client
	Service        *spreadsheet.Service
	sheetsId       string
	Spreadsheet    spreadsheet.Spreadsheet
	Sheet          *spreadsheet.Sheet
	sheetIndex     uint
	sheetTitle     string
	KeyColumnIndex uint
	Columns        []Column
	ColumnMapKeyLc map[string]Column
	ItemMap        map[string]Item
}

func (sm *SheetsMap) SheetTitle() string {
	return sm.sheetTitle
}

func (sm *SheetsMap) ColumnsKeys() []string {
	keys := []string{}
	for _, col := range sm.Columns {
		keys = append(keys, col.Name)
	}
	return keys
}

func (sm *SheetsMap) DataColumnsKeys() []string {
	keys := []string{}
	for i, col := range sm.Columns {
		if i < 2 {
			continue
		}
		keys = append(keys, col.Name)
	}
	return keys
}

func NewSheetsMap() SheetsMap {
	return SheetsMap{
		GoogleClient:   nil,
		Service:        nil,
		sheetsId:       "",
		Columns:        []Column{},
		ColumnMapKeyLc: map[string]Column{},
		ItemMap:        map[string]Item{},
	}
}

func NewSheetsMapIndex(googleClient *http.Client, spreadsheetId string, sheetIndex uint) (SheetsMap, error) {
	sm := NewSheetsMap()
	sm.GoogleClient = googleClient
	sm.Service = spreadsheet.NewServiceWithClient(googleClient)

	spreadsheet, err := sm.Service.FetchSpreadsheet(spreadsheetId)
	if err != nil {
		return sm, err
	}
	sm.Spreadsheet = spreadsheet

	sheet, err := sm.Spreadsheet.SheetByIndex(sheetIndex)
	if err != nil {
		return sm, err
	}
	sm.sheetIndex = sheetIndex
	sm.Sheet = sheet

	return sm, nil
}

func NewSheetsMapTitle(googleClient *http.Client, spreadsheetId string, sheetTitle string) (SheetsMap, error) {
	sm := NewSheetsMap()
	sm.GoogleClient = googleClient
	sm.Service = spreadsheet.NewServiceWithClient(googleClient)

	spreadsheet, err := sm.Service.FetchSpreadsheet(spreadsheetId)
	if err != nil {
		return sm, err
	}
	sm.Spreadsheet = spreadsheet

	sheet, err := sm.Spreadsheet.SheetByTitle(sheetTitle)
	if err != nil {
		return sm, err
	}
	sm.sheetTitle = sheetTitle
	sm.Sheet = sheet

	return sm, nil
}

type Item struct {
	Key     string
	Display string
	Row     uint
	Data    map[string]string
}

func (item *Item) ItemDisplayOrKey() string {
	display := strings.TrimSpace(item.Display)
	if len(display) > 0 {
		return display
	}
	return strings.TrimSpace(item.Key)
}

type Column struct {
	Name               string
	NameAliases        []string
	Abbreviation       string
	Index              uint64
	Enums              []Enum
	AliasLcToCanonical map[string]string
	InfoURLs           []InfoURL
}

type InfoURL struct {
	Text string
	URL  string
}

func NewColumn() Column {
	return Column{
		NameAliases:        []string{},
		Enums:              []Enum{},
		AliasLcToCanonical: map[string]string{},
		InfoURLs:           []InfoURL{}}
}

func (col *Column) AddEnum(enum Enum) {
	col.Enums = append(col.Enums, enum)
	col.AliasLcToCanonical[strings.ToLower(enum.Canonical)] = enum.Canonical
	for _, alias := range enum.Aliases {
		col.AliasLcToCanonical[strings.ToLower(alias)] = enum.Canonical
	}
}

func (col *Column) EnumsCanonical() []string {
	canonicals := []string{}
	for _, enum := range col.Enums {
		canonicals = append(canonicals, enum.Canonical)
	}
	return canonicals
}

func (col *Column) EnumsStrings() []string {
	enums := []string{}
	for _, enum := range col.Enums {
		enums = append(enums, enum.Values()...)
	}
	return enums
}

type Enum struct {
	Canonical string
	Aliases   []string
}

func (enum *Enum) Values() []string {
	values := []string{}
	if len(enum.Canonical) > 0 {
		values = append(values, enum.Canonical)
	}
	if len(enum.Aliases) > 0 {
		values = append(values, enum.Aliases...)
	}
	return values
}

// ParseColumn
// tshirt size - XS, S, M, L, XL, XXL, XXXL
// colName | colAbbr | Enums | URLs
func ParseColumn(input string) (Column, error) {
	parts := strings.Split(input, " - ")
	col := NewColumn()
	if len(parts) <= 0 {
		return col, fmt.Errorf("Column Format Error for [%v]", input)
	} else if len(parts) > 4 {
		return col, fmt.Errorf("Column Format Error for [%v]", input)
	}

	if len(parts) >= 1 {
		colNames := stringsutil.SplitCondenseSpace(parts[0], "|")
		if len(colNames) > 0 {
			col.Name = colNames[0]
		}
		if len(colNames) > 1 {
			for i := 1; i < len(colNames); i++ {
				col.NameAliases = append(col.NameAliases, colNames[i])
			}
		}
		// col.Value = strings.TrimSpace(parts[0])
	}
	if len(parts) >= 2 { // Have Abbrevations
		col.Abbreviation = strings.TrimSpace(parts[1])
	}

	if len(parts) >= 3 { // Have enum values
		enums := stringsutil.SplitCondenseSpace(parts[2], ",")
		for _, enumPlus := range enums {
			enumVariations := stringsutil.SplitCondenseSpace(enumPlus, "|")
			if len(enumVariations) > 0 {
				enum := Enum{Canonical: enumVariations[0]}
				if len(enumVariations) > 1 { // Have aliases
					enum.Aliases = enumVariations[1:]
				}
				col.AddEnum(enum)
			}
		}
	}
	if len(parts) >= 4 { // URL
		urls := strings.Split(parts[3], " ~ ")
		for _, urlInfoRaw := range urls {
			urlInfoParts := stringsutil.SplitCondenseSpace(urlInfoRaw, "|")
			if len(urlInfoParts) == 2 {
				col.InfoURLs = append(col.InfoURLs, InfoURL{
					Text: urlInfoParts[0],
					URL:  urlInfoParts[1],
				})
			}
		}
	}
	return col, nil
}

func (col *Column) ValueToCanonical(val string) (string, error) {
	if len(col.Enums) == 0 {
		return val, nil
	}
	valLc := TrimSpaceToLower(val)

	for tryLc, tryCanonical := range col.AliasLcToCanonical {
		tryLc = TrimSpaceToLower(tryLc)
		if valLc == tryLc {
			return strings.TrimSpace(tryCanonical), nil
		}
	}

	enums := strings.Join(col.EnumsCanonical(), ", ")
	return enums, fmt.Errorf("Column [%v] Value [%v] not valid [%v]", col.Name, val, enums)
}

func (sm *SheetsMap) FullRead() error {
	err := sm.ReadColumns()
	if err != nil {
		return err
	}
	return sm.ReadItems()
}

func (sm *SheetsMap) ReadColumns() error {
	colsMap := map[string]Column{}
	colsArr := []Column{}

	for _, row := range sm.Sheet.Rows {
		for j, cell := range row {
			colValRaw := strings.TrimSpace(cell.Value)
			if len(colValRaw) < 1 {
				break
			}

			col, err := ParseColumn(colValRaw)
			if err != nil {
				return err
			}
			col.Index = uint64(j)
			colKeyParsedLc := strings.ToLower(col.Name)

			colsArr = append(colsArr, col)
			if _, ok := colsMap[colKeyParsedLc]; ok {
				return fmt.Errorf("Duplicate column names for: %v", colValRaw)
			}
			colsMap[colKeyParsedLc] = col
		}
		break
	}

	sm.Columns = colsArr
	sm.ColumnMapKeyLc = colsMap
	return nil
}

func (sm *SheetsMap) ReadItems() error {
	itemMap := map[string]Item{}

	for i, row := range sm.Sheet.Rows {
		if i == 0 {
			continue
		}
		item := Item{
			Row:  uint(i),
			Data: map[string]string{},
		}
		for j, cell := range row {
			if j >= len(sm.Columns) {
				break
			}

			val := cell.Value
			if j == 0 {
				item.Key = val
			} else if j == 1 {
				item.Display = val
			}
			col := sm.Columns[j]
			item.Data[col.Name] = val
		}
		if _, ok := itemMap[item.Key]; ok {
			return fmt.Errorf("Duplicate key names for: %v", item.Key)
		}
		itemMap[item.Key] = item
	}

	sm.ItemMap = itemMap
	return nil
}

func (sm *SheetsMap) GetItem(key string) (Item, error) {
	if item, ok := sm.ItemMap[key]; !ok {
		return Item{}, fmt.Errorf("Cannot find key %v", key)
	} else {
		return item, nil
	}
}

func (sm *SheetsMap) GetItemProperty(key string, val string) (string, error) {
	item, err := sm.GetItem(key)
	if err != nil {
		return "", err
	}
	val, ok := item.Data[val]
	if !ok {
		return "", fmt.Errorf("Cannot find value for property [%v]", val)
	}
	return val, nil
}

func (sm *SheetsMap) IsItemComplete(item *Item) bool {
	complete := true
	for _, col := range sm.Columns {
		if val, ok := item.Data[col.Name]; !ok || len(strings.TrimSpace(val)) == 0 {
			complete = false
			break
		}
	}
	return complete
}

func (sm *SheetsMap) IsItemPartial(item *Item) bool {
	completeCols := 0
	for _, col := range sm.Columns {
		if val, ok := item.Data[col.Name]; ok && len(strings.TrimSpace(val)) > 0 {
			completeCols += 1
		}
	}
	switch completeCols {
	case len(sm.Columns):
		return false
	case 0:
		return false
	default:
		return true
	}
}

func (sm *SheetsMap) GetOrCreateItemWithName(itemKey, itemName string) (Item, error) {
	itemKey = TrimSpaceToLower(itemKey)
	itemName = strings.TrimSpace(itemName)
	if item, ok := sm.ItemMap[itemKey]; !ok {
		item := Item{
			Key:     itemKey,
			Display: itemName,
			Data:    map[string]string{},
		}
		if len(sm.Columns) > 0 {
			item.Data[sm.Columns[0].Name] = itemKey
		}
		if len(sm.Columns) > 1 {
			item.Data[sm.Columns[1].Name] = itemName
		}

		itemCount := len(sm.ItemMap)
		nextRowIdx := itemCount + 1
		item.Row = uint(nextRowIdx)

		sm.Sheet.Update(nextRowIdx, 0, itemKey)
		sm.Sheet.Update(nextRowIdx, 1, itemName)
		err := sm.Sheet.Synchronize()
		if err == nil {
			sm.ItemMap[itemKey] = item
		}
		return item, err
	} else {
		if item.Display != itemName {
			item.Display = itemName
			if len(sm.Columns) > 1 {
				item.Data[sm.Columns[1].Name] = itemName
			}
			err := sm.SynchronizeItem(item)
			return item, err
		}
		return item, nil
	}
}

func (sm *SheetsMap) GetOrCreateItem(itemKey string) (Item, error) {
	itemKey = TrimSpaceToLower(itemKey)
	if item, ok := sm.ItemMap[itemKey]; !ok {
		item := Item{
			Key:  itemKey,
			Data: map[string]string{},
		}
		if len(sm.Columns) > 0 {
			item.Data[sm.Columns[0].Name] = itemKey
		}

		itemCount := len(sm.ItemMap)
		nextRowIdx := itemCount + 1
		item.Row = uint(nextRowIdx)

		sm.Sheet.Update(nextRowIdx, 0, itemKey)
		err := sm.Sheet.Synchronize()
		if err == nil {
			sm.ItemMap[itemKey] = item
		}
		return item, err
	} else {
		return item, nil
	}
}

func (sm *SheetsMap) UpdateItem(item Item, key, val string, synchronize bool) (string, error) {
	// Get key column
	keyLc := TrimSpaceToLower(key)
	col, ok := sm.ColumnMapKeyLc[keyLc]
	if !ok {
		return "", fmt.Errorf("Key Not Found: %v", key)
	}

	// Process value
	str, err := col.ValueToCanonical(val)
	if err != nil {
		return str, err
	}

	item.Data[keyLc] = str
	return "", sm.SynchronizeItem(item)
}

func (sm *SheetsMap) SynchronizeItem(item Item) error {
	rowIdx := item.Row
	for colIdx, col := range sm.Columns {
		if val, ok := item.Data[col.Name]; ok {
			sm.Sheet.Update(int(rowIdx), colIdx, val)
		} else {
			sm.Sheet.Update(int(rowIdx), colIdx, "")
		}
	}
	sm.Sheet.Update(int(rowIdx), 1, item.Display)
	return sm.Sheet.Synchronize()
}

func (sm *SheetsMap) EmptyCols(item Item) []string {
	//rowIdx := item.Row
	emptyCols := []string{}
	for _, col := range sm.Columns {
		if val, ok := item.Data[col.Name]; ok {
			if len(strings.TrimSpace(val)) == 0 {
				emptyCols = append(emptyCols, col.Name)
			}
		} else {
			emptyCols = append(emptyCols, col.Name)
		}
	}
	return emptyCols
}

func (sm *SheetsMap) SetItemKeyColValue(itemKey, colKeyRaw, colValRaw string) (Item, error) {
	item, err := sm.GetOrCreateItem(itemKey)
	if err != nil {
		return item, err
	}

	colKey := strings.TrimSpace(colKeyRaw)
	colKeyLc := strings.ToLower(colKeyRaw)
	col, ok := sm.ColumnMapKeyLc[colKeyLc]
	if !ok {
		return item, fmt.Errorf("%s [%s]", ErrorColumnNotFound, colKey)
	}

	colVal, err := col.ValueToCanonical(colValRaw)
	if err != nil {
		//enumsCanonical := []string{}
		//fmt.Errorf("%s [%s] [%s]", ErrorEnumNotMatched, colValRaw, strings.Join(col.EnumsCanonical(), ", "))
		return item, err
	}

	item.Data[colKey] = colVal
	sm.ItemMap[itemKey] = item
	return item, nil
}

type Intent struct {
	Name  string
	Slots map[string]string
}

func TrimSpaceToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type Stat struct {
	Name  string
	Names []string
	Count int64
}

func (sm *SheetsMap) CombinedStatsCol0Enum() ([]Stat, error) {
	permutationsMap := map[string]map[string]int64{}
	col0 := Column{}
	for _, item := range sm.ItemMap {
		vals := []string{}
		for i, col := range sm.Columns {
			if i < 2 {
				continue
			}
			if i == 2 {
				col0 = col
			}
			if colVal, ok := item.Data[col.Name]; ok {
				colVal = strings.TrimSpace(colVal)
				if len(colVal) > 0 {
					vals = append(vals, colVal)
				} else {
					vals = append(vals, "?")
				}
			} else {
				vals = append(vals, "?")
			}
		}
		if len(vals) > 0 {
			val0 := vals[0]
			valsStr := strings.Join(vals, ", ")
			if _, ok := permutationsMap[val0]; !ok {
				permutationsMap[val0] = map[string]int64{}
			}
			if _, ok := permutationsMap[val0][valsStr]; !ok {
				permutationsMap[val0][valsStr] = 1
			} else {
				permutationsMap[val0][valsStr] += 1
			}
		}

	}

	stats := []Stat{}

	for _, enum := range col0.Enums {
		if mss, ok := permutationsMap[enum.Canonical]; ok {
			valsStrs := []string{}
			for valsStr := range mss {
				valsStrs = append(valsStrs, valsStr)
			}
			sort.Strings(valsStrs)
			for _, valsStr := range valsStrs {
				keyCount := mss[valsStr]
				stats = append(stats, Stat{
					Name:  valsStr,
					Count: int64(keyCount),
				})
			}
		}
	}
	enumCanonicalUnknown := "?"
	if mss, ok := permutationsMap[enumCanonicalUnknown]; ok {
		valsStrs := []string{}
		for valsStr := range mss {
			valsStrs = append(valsStrs, valsStr)
		}
		sort.Strings(valsStrs)
		for _, valsStr := range valsStrs {
			keyCount := mss[valsStr]
			stats = append(stats, Stat{
				Name:  valsStr,
				Count: int64(keyCount),
			})
		}
	}
	return stats, nil
}

func (sm *SheetsMap) SetItemKeyDisplay(itemKey, itemDisplay string) error {
	item, err := sm.GetOrCreateItem(itemKey)
	if err != nil {
		return err
	}
	itemDisplay = strings.TrimSpace(itemDisplay)
	if item.Display != itemDisplay {
		item.Display = itemDisplay
		return sm.SynchronizeItem(item)
	}
	return nil
}

func (sm *SheetsMap) SetItemKeyString(itemKey, cmdRaw string) (Intent, error) {
	cmdRawLc := TrimSpaceToLower(cmdRaw)
	intent := Intent{Slots: map[string]string{}}

	for _, col := range sm.ColumnMapKeyLc {
		numColNameAll := 1 + len(col.NameAliases)
		for i := 0; i < numColNameAll; i++ {
			colNameTry := ""
			if i == 0 {
				colNameTry = col.Name
			} else {
				colNameTry = col.NameAliases[i-1]
			}

			colNameTryLc := TrimSpaceToLower(colNameTry)
			//pat := fmt.Sprintf("^%v\\s*(.*)$", colNameTryLc)
			pat := `^` + colNameTryLc + `\s*(.*)$`
			pat1, err := regexp.Compile(pat)
			if err != nil {
				return intent, err
			}
			m := pat1.FindStringSubmatch(cmdRawLc)
			if len(m) == 2 {
				valCanonical, err := col.ValueToCanonical(m[1])
				if err != nil {
					return intent, err
				}
				item, err := sm.SetItemKeyColValue(itemKey, col.Name, valCanonical)
				if err != nil {
					return intent, err
				}
				err = sm.SynchronizeItem(item)
				if err != nil {
					return intent, err
				}
				return intent, nil
			}
		}
	}

	return intent, fmt.Errorf("E_CANNOT_FIND_MATCH KEY[%v] CMD[%v]", itemKey, cmdRaw)
}
