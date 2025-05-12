package ids

import (
	"kiokun-go/dictionaries/common"
)

func init() {
	// Register the basic UCS IDS file
	common.RegisterDictionary("ids", "IDS-UCS-Basic.txt", &Importer{})

	// Register the Extension A IDS file as a separate dictionary
	common.RegisterDictionary("ids_ext_a", "IDS-UCS-Ext-A.txt", &Importer{})

	// Additional extension files can be registered here as needed
}
