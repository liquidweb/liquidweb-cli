package apiTypes

type PaginatedList struct {
	ItemCount int64                    `json:"item_count" mapstructure:"item_count"`
	ItemTotal int64                    `json:"item_total" mapstructure:"item_total"`
	Items     []map[string]interface{} `json:"items" mapstructure:"items"`
	PageNum   int64                    `json:"page_num" mapstructure:"page_num"`
	PageSize  int64                    `json:"page_size mapstructure:"page_size"`
	PageTotal int64                    `json:"page_total" mapstructure:"page_total"`
}

type MergedPaginatedList struct {
	Items       []map[string]interface{} `json:"items" mapstructure:"items"`
	MergedPages int64                    `json:"merged_pages" mapstructure:"merged_pages"`
	PageSize    int64                    `json:"page_size" mapstructure:"page_size"`
}
