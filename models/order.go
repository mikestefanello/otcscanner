package models

import (
	"errors"
	"fmt"
	"strconv"
)

// Order describes an order
type Order struct {
	PackageID                              string `bson:"packageId" csv:"Package ID" validate:"required"`
	SenderFirstName                        string `bson:"senderFirstName" csv:"Sender First Name"`
	SenderLastName                         string `bson:"senderLastName" csv:"Sender Last Name"`
	SenderBusinessName                     string `bson:"senderBusinessName" csv:"Sender Business Name"`
	SenderAddressLine1                     string `bson:"senderAddressLine1" csv:"Sender Address Line 1"`
	SenderAddressLine2                     string `bson:"senderAddressLine2" csv:"Sender Address Line 2"`
	SenderCity                             string `bson:"senderCity" csv:"Sender City"`
	SenderProvince                         string `bson:"senderProvince" csv:"Sender Province"`
	SenderPostalCode                       string `bson:"senderPostalCode" csv:"Sender Postal Code"`
	SenderCountryCode                      string `bson:"senderCountryCode" csv:"Sender Country Code"`
	SenderPhoneNumber                      string `bson:"senderPhoneNumber" csv:"Sender Phone Number"`
	RecipientFirstName                     string `bson:"recipientFirstName" csv:"Recipient First Name"`
	RecipientLastName                      string `bson:"recipientLastName" csv:"Recipient Last Name"`
	RecipientBusinessName                  string `bson:"recipientBusinessName" csv:"Recipient Business Name"`
	RecipientAddressLine1                  string `bson:"recipientAddressLine1" csv:"Recipient Address Line 1"`
	RecipientAddressLine2                  string `bson:"recipientAddressLine2" csv:"Recipient Address Line 2"`
	RecipientAddressLine3                  string `bson:"recipientAddressLine3" csv:"Recipient Address Line 3"`
	RecipientInLineTranslationAddressLine1 string `bson:"recipientInLineTranslationAddressLine1" csv:"RecipientInLineTranslationAddressLine1"`
	RecipientInLineTranslationAddressLine2 string `bson:"recipientInLineTranslationAddressLine2" csv:"RecipientInLineTranslationAddressLine2"`
	RecipientCity                          string `bson:"recipientCity" csv:"Recipient City"`
	RecipientProvince                      string `bson:"recipientProvince" csv:"Recipient Province"`
	RecipientPostalCode                    string `bson:"recipientPostalCode" csv:"Recipient Postal Code"`
	RecipientCountryCode                   string `bson:"recipientCountryCode" csv:"Recipient Country Code"`
	RecipientPhoneNumber                   string `bson:"recipientPhoneNumber" csv:"Recipient Phone Number"`
	RecipientEmailAddress                  string `bson:"recipientEmailAddress" csv:"Recipient E-mail Address"`
	PackageWeight                          string `bson:"packageWeight" csv:"Package Weight"`
	WeightUnit                             string `bson:"weightUnit" csv:"Weight Unit"`
	ServiceType                            string `bson:"serviceType" csv:"Service Type"`
	RateType                               string `bson:"rateType" csv:"Rate Type"`
	PackageType                            string `bson:"packageType" csv:"Package Type"`
	PackagePhysicalCount                   string `bson:"packagePhysicalCount" csv:"Package Physical Count"`
	PFCEELCode                             string `bson:"pfcEelCode" csv:"PFC/EEL Code"`
	ItemID                                 string `bson:"itemId" csv:"Item ID"`
	ItemDescription                        string `bson:"itemDescription" csv:"Item Description"`
	UnitValueUSD                           string `bson:"unitValueUsd" csv:"Unit Value (USD)"`
	Quantity                               string `bson:"quantity" csv:"Quantity"`
	CountryOfOrigin                        string `bson:"countryOfOrigin" csv:"Country Of Origin"`
	Country                                string `bson:"country" csv:"Country"`
	Weight                                 string `bson:"weight" csv:"Weight"`
	Service                                string `bson:"service" csv:"Service"`
	Length                                 string `bson:"length" csv:"Length"`
	Width                                  string `bson:"width" csv:"Width"`
	Height                                 string `bson:"height" csv:"Height"`
	DIM                                    string `bson:"dim" csv:"DIM"`
	Account                                string `bson:"account" csv:"Account"`
}

// Orders is a slice of order structs
type Orders []Order

// CalculateDim calculates and sets the DIM field on a given order
func (o *Order) CalculateDim() error {
	// Check if all dimensions are populated
	if o.Length != "" && o.Width != "" && o.Height != "" {
		length, err := strconv.ParseFloat(o.Length, 64)
		if err != nil {
			return errors.New("Unable to parse length")
		}

		width, err := strconv.ParseFloat(o.Width, 64)
		if err != nil {
			return errors.New("Unable to parse width")
		}

		height, err := strconv.ParseFloat(o.Height, 64)
		if err != nil {
			return errors.New("Unable to parse height")
		}

		dim := fmt.Sprintf("%.2f", (length*width*height)/139)
		o.DIM = dim
	}

	return nil
}
