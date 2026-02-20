package models

type ChartSeries struct {
    Label     string  `json:"label" gorm:"column:label"`           // Mapping ke kolom 'label' di SQL
    Target    float64 `json:"target" gorm:"column:target"`         // Mapping ke kolom 'target'
    Actual    float64 `json:"actual" gorm:"column:actual"`         // Mapping ke kolom 'actual'
    ActualOK  float64 `json:"actual_ok,omitempty" gorm:"column:actual_ok"` 
    ActualNG  float64 `json:"actual_ng,omitempty" gorm:"column:actual_ng"`
    
    // Ini kolom barumu
    ItemName  string  `json:"item_name" gorm:"column:item_name"`   
    
    ExtraInfo string  `json:"extra_info,omitempty" gorm:"column:extra_info"` // Penting: mapping ke snake_case SQL
}