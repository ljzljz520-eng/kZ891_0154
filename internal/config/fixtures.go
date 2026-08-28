package config

import "warrantyservice/internal/model"

func FixtureWarranties() []model.WarrantyRecord {
	return []model.WarrantyRecord{
		{ID: 1, Phone: "13800138000", SerialNumber: "TV-2025-0001", ExpiryDate: "2027-05-31", ServicePoints: []string{"上海徐汇服务中心", "上海浦东服务中心"}, PurchaseChannel: "官方商城", Category: "电视", UpdatedAt: "2025-01-10", Active: true},
		{ID: 2, Phone: "13900139000", SerialNumber: "WM-2024-0042", ExpiryDate: "2025-05-30", ServicePoints: []string{"北京朝阳服务中心"}, PurchaseChannel: "授权门店", Category: "洗衣机", UpdatedAt: "2024-06-12", Active: true},
		{ID: 3, Phone: "13700137000", SerialNumber: "AC-2025-0088", ExpiryDate: "2028-02-28", ServicePoints: []string{"广州天河服务中心"}, PurchaseChannel: "电商旗舰店", Category: "空调", UpdatedAt: "2025-02-01", Active: true},
	}
}

func FixturePoints() []model.ServicePoint {
	return []model.ServicePoint{
		{ID: 1, Name: "上海徐汇服务中心", Address: "上海市徐汇区漕溪北路88号", Phone: "021-64001234", Categories: []string{"电视", "冰箱"}, OpenHours: "09:00-18:00"},
		{ID: 2, Name: "上海浦东服务中心", Address: "上海市浦东新区张江路100号", Phone: "021-58005678", Categories: []string{"电视", "空调"}, OpenHours: "09:00-18:00"},
		{ID: 3, Name: "北京朝阳服务中心", Address: "北京市朝阳区望京街12号", Phone: "010-64007890", Categories: []string{"洗衣机", "冰箱"}, OpenHours: "08:30-17:30"},
		{ID: 4, Name: "广州天河服务中心", Address: "广州市天河区体育西路66号", Phone: "020-38001234", Categories: []string{"空调", "电视"}, OpenHours: "09:00-18:00"},
	}
}
