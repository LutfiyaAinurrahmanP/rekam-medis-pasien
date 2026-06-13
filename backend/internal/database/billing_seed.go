package database

import (
	"fmt"
	"log"
	"time"

	"github.com/LutfiyaAinurrahmanP/sirekam-medis-pasien/internal/models"
	"gorm.io/gorm"
)

func SeedBilling(db *gorm.DB) error {
	var count int64
	db.Unscoped().Model(&models.Billing{}).Count(&count)
	if count > 0 {
		log.Println("Billing table is not empty, skipping seed")
		return nil
	}

	log.Println("Seeding Billing data...")

	var billings []models.Billing
	now := time.Now()

	// 1. Generate 22 active billings
	for i := 1; i <= 22; i++ {
		patientID := uint((i % 5) + 1)
		var medicalRecordID *uint
		if i%2 == 0 {
			mID := uint((i % 5) + 1)
			medicalRecordID = &mID
		}

		status := "pending"
		var paymentMethod *string
		paidAmount := 0.0

		if i%3 == 0 {
			status = "paid"
			method := "cash"
			if i%2 == 0 {
				method = "transfer"
			}
			paymentMethod = &method
			paidAmount = 500000.0
		} else if i%4 == 0 {
			status = "partial"
			paidAmount = 250000.0
		} else if i%5 == 0 {
			status = "cancelled"
		}

		totalAmount := 500000.0
		remainingAmount := totalAmount - paidAmount

		methodVal := paymentMethod
		notes := fmt.Sprintf("Catatan billing %d", i)

		billing := models.Billing{
			PatientID:       patientID,
			MedicalRecordID: medicalRecordID,
			InvoiceNumber:   fmt.Sprintf("INV-2024-%06d", i),
			BillingDate:     now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			DueDate:         now.AddDate(0, 0, 7-(i%10)).Format("2006-01-02"),
			TotalAmount:     totalAmount,
			PaidAmount:      paidAmount,
			RemainingAmount: remainingAmount,
			Status:          status,
			PaymentMethod:   methodVal,
			Notes:           &notes,
			CreatedAt:       now,
			UpdatedAt:       now,
			Items: []models.BillingItem{
				{
					Description: "Biaya Konsultasi",
					Quantity:        1,
					UnitPrice:       150000,
					TotalPrice:      150000,
				},
				{
					Description: "Obat Paracetamol",
					Quantity:        10,
					UnitPrice:       5000,
					TotalPrice:      50000,
				},
			},
		}
		
		billings = append(billings, billing)
	}

	// 2. Generate deleted billings
	for i := 23; i <= 44; i++ {
		patientID := uint((i % 5) + 1)

		billings = append(billings, models.Billing{
			PatientID:     patientID,
			InvoiceNumber: fmt.Sprintf("INV-DEL-%06d", i),
			BillingDate:     now.AddDate(0, 0, -(i % 10)).Format("2006-01-02"),
			DueDate:         now.AddDate(0, 0, 7-(i%10)).Format("2006-01-02"),
			TotalAmount:   100000,
			PaidAmount:    0,
			RemainingAmount: 100000,
			Status: "pending",
			CreatedAt:     now,
			UpdatedAt:     now,
			DeletedAt:     gorm.DeletedAt{Time: now, Valid: true},
		})
	}

	if err := db.Create(&billings).Error; err != nil {
		log.Printf("Error seeding Billing: %v\n", err)
		return err
	}

	// update items type
	var items []models.BillingItem
	db.Find(&items)
	for _, item := range items {
		typ := "consultation"
		if item.Description == "Obat Paracetamol" {
			typ = "medicine"
		}
		item.ItemType = &typ
		db.Save(&item)
	}

	log.Println("Billing seeded successfully")
	return nil
}
