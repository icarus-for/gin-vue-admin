package service

import (
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/dto"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

type Address struct{}

func (g *Address) GetAddressListInfo(userId string) ([]dto.ShippingAddressInfo, error) {
	var addressList []dto.ShippingAddressInfo
	err := global.GVA_DB.Table("user_shipping_addresses").
		Select(`
			id AS address_id, postal_code, prefecture, city, address_line1, address_line2, recipient_name, phone_number,
			CASE WHEN is_default = 1 THEN 'true' ELSE 'false' END AS is_default
			`).
		Where("user_id = ?", userId).
		Order("is_default DESC, updated_at DESC").
		Scan(&addressList).Error
	if err != nil {
		return nil, err
	}
	return addressList, err
}

func (g *Address) AddAddressListInfo(userId string, postalCode string, Prefecture string, City string, AddressLine1 string, AddressLine2 string, RecipientName string, PhoneNumber string, IsDefault int) (string, error) {
	if IsDefault == 1 {
		err := global.GVA_DB.Table("user_shipping_addresses").
			Where("user_id= ? ", userId).
			Update("is_default", 0).Error
		if err != nil {
			return "", err
		}
	}
	err := global.GVA_DB.Table("user_shipping_addresses").
		Create(map[string]interface{}{
			"user_id":        userId,
			"postal_code":    postalCode,
			"prefecture":     Prefecture,
			"city":           City,
			"address_line1":  AddressLine1,
			"address_line2":  AddressLine2,
			"recipient_name": RecipientName,
			"phone_number":   PhoneNumber,
			"is_default":     IsDefault,
		}).Error
	if err != nil {
		return "", err
	}
	return "アドレスの追加は成功しました。", err
}

func (g *Address) DelAddressListInfo(userId string, addressId string) (string, error) {
	var addressCount int64
	err := global.GVA_DB.Table("user_shipping_addresses").
		Select("id,user_id").
		Where("id = ? AND user_id=?", addressId, userId).
		Count(&addressCount).Error
	if err != nil {
		return "", err
	}
	if addressCount == 0 {
		return "", fmt.Errorf("住址id不存在,无法删除")
	}
	err = global.GVA_DB.Table("user_shipping_addresses").
		Where("id = ? AND user_id=?", addressId, userId).
		Delete(nil).Error
	if err != nil {
		return "", err
	}
	return "配送住所から削除しました。", err
}

func (g *Address) ChangeAddressInfo(userId string, addressId string, postalCode string, Prefecture string, City string, AddressLine1 string, AddressLine2 string, RecipientName string, PhoneNumber string, IsDefault int) (string, error) {
	var addressCount int64
	err := global.GVA_DB.Table("user_shipping_addresses").
		Select("id,user_id").
		Where("id = ? AND user_id=?", addressId, userId).
		Count(&addressCount).Error
	if err != nil {
		return "", err
	}
	if addressCount == 0 {
		return "", fmt.Errorf("住址id不存在,无法更改")
	}
	updateData := map[string]interface{}{}
	if postalCode != "" {
		updateData["postal_code"] = postalCode
	}
	if Prefecture != "" {
		updateData["prefecture"] = Prefecture
	}
	if City != "" {
		updateData["city"] = City
	}
	if AddressLine1 != "" {
		updateData["address_line1"] = AddressLine1
	}
	if AddressLine2 != "" {
		updateData["address_line2"] = AddressLine2
	}
	if RecipientName != "" {
		updateData["recipient_name"] = RecipientName
	}
	if PhoneNumber != "" {
		updateData["phone_number"] = PhoneNumber
	}
	if IsDefault != 0 {
		if IsDefault == 1 {
			err := global.GVA_DB.Table("user_shipping_addresses").
				Where("user_id= ? ", userId).
				Update("is_default", 0).Error
			if err != nil {
				return "", err
			}
		}
		updateData["recipient_name"] = RecipientName
	}
	if len(updateData) > 0 {
		err := global.GVA_DB.Table("user_shipping_addresses").
			Where("user_id=? AND id=?", userId, addressId).
			Updates(updateData).Error
		if err != nil {
			return "", err
		}
	}

	return "配送住址を更新しました。", nil
}
