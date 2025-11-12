package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"evermos-project/models"
	"evermos-project/utils"

	"github.com/gofiber/fiber/v2"
)

type ProvCityController struct{}

func NewProvCityController() *ProvCityController {
	return &ProvCityController{}
}

const (
	provincesURL = "https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json"
	provinceURL  = "https://www.emsifa.com/api-wilayah-indonesia/api/province/%s.json"
	citiesURL    = "https://www.emsifa.com/api-wilayah-indonesia/api/regencies/%s.json"
	cityURL      = "https://www.emsifa.com/api-wilayah-indonesia/api/regency/%s.json"
	districtsURL = "https://www.emsifa.com/api-wilayah-indonesia/api/districts/%s.json"
	districtURL  = "https://www.emsifa.com/api-wilayah-indonesia/api/district/%s.json"
	villagesURL  = "https://www.emsifa.com/api-wilayah-indonesia/api/villages/%s.json"
	villageURL   = "https://www.emsifa.com/api-wilayah-indonesia/api/village/%s.json"
)

// GET ALL PROVINCES
func (pcc *ProvCityController) GetListProvinces(c *fiber.Ctx) error {
	resp, err := http.Get(provincesURL)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var provinces []models.Province
	if err := json.NewDecoder(resp.Body).Decode(&provinces); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", provinces)
}

// GET LIST CITIES BY PROVINCE
func (pcc *ProvCityController) GetListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	url := fmt.Sprintf(citiesURL, provID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var cities []models.City
	if err := json.NewDecoder(resp.Body).Decode(&cities); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", cities)
}

// GET DETAIL PROVINCE
func (pcc *ProvCityController) GetDetailProvince(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	url := fmt.Sprintf(provinceURL, provID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var province models.Province
	if err := json.NewDecoder(resp.Body).Decode(&province); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", province)
}

// GET DETAIL CITY
func (pcc *ProvCityController) GetDetailCity(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	url := fmt.Sprintf(cityURL, cityID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var city models.City
	if err := json.NewDecoder(resp.Body).Decode(&city); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", city)
}

// GET LIST DISTRICTS BY CITY/REGENCY
func (pcc *ProvCityController) GetListDistricts(c *fiber.Ctx) error {
	cityID := c.Params("city_id")
	url := fmt.Sprintf(districtsURL, cityID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var districts []models.District
	if err := json.NewDecoder(resp.Body).Decode(&districts); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", districts)
}

// GET DETAIL DISTRICT
func (pcc *ProvCityController) GetDetailDistrict(c *fiber.Ctx) error {
	districtID := c.Params("district_id")
	url := fmt.Sprintf(districtURL, districtID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var district models.District
	if err := json.NewDecoder(resp.Body).Decode(&district); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", district)
}

// GET LIST VILLAGES BY DISTRICT
func (pcc *ProvCityController) GetListVillages(c *fiber.Ctx) error {
	districtID := c.Params("district_id")
	url := fmt.Sprintf(villagesURL, districtID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var villages []models.Village
	if err := json.NewDecoder(resp.Body).Decode(&villages); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", villages)
}

// GET DETAIL VILLAGE
func (pcc *ProvCityController) GetDetailVillage(c *fiber.Ctx) error {
	villageID := c.Params("village_id")
	url := fmt.Sprintf(villageURL, villageID)

	resp, err := http.Get(url)
	if err != nil {
		return utils.FiberError(c, "Failed to get data", err)
	}
	defer resp.Body.Close()

	var village models.Village
	if err := json.NewDecoder(resp.Body).Decode(&village); err != nil {
		return utils.FiberError(c, "Failed to decode JSON", err)
	}

	return utils.FiberSuccess(c, "Succeed to get data", village)
}