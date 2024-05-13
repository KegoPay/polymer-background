package toexcel

import (
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
	"usepolymer.co/background/logger"
	"usepolymer.co/background/models"
	"usepolymer.co/background/utils"
)

func TransactionToExcel(transactions *[]models.Transaction) (*string, error) {
	file := excelize.NewFile()
	file.SetCellValue("Sheet1", "A1", "Date")
	file.SetCellValue("Sheet1", "B1", "Transaction ID")
	file.SetCellValue("Sheet1", "C1", "Narration")
	file.SetCellValue("Sheet1", "D1", "Amount")
	file.SetCellValue("Sheet1", "E1", "Currency")
	file.SetCellValue("Sheet1", "F1", "Intent")
	file.SetCellValue("Sheet1", "G1", "Recipient Name")
	file.SetCellValue("Sheet1", "H1", "Recipient Account Number")
	file.SetCellValue("Sheet1", "I1", "Recipient Bank")

	for i, transaction := range *transactions {
		file.SetCellValue("Sheet1", fmt.Sprintf("A%d", i+2), transaction.CreatedAt.Format("2006-01-02 15:04:05"))
		file.SetCellValue("Sheet1", fmt.Sprintf("B%d", i+2), transaction.ID)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%d", i+2), transaction.Description)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%d", i+2), transaction.Amount)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%d", i+2), transaction.Currency)
		file.SetCellValue("Sheet1", fmt.Sprintf("F%d", i+2), transaction.Intent)
		file.SetCellValue("Sheet1", fmt.Sprintf("G%d", i+2), transaction.Recepient.FullName)
		file.SetCellValue("Sheet1", fmt.Sprintf("H%d", i+2), transaction.Recepient.AccountNumber)
		file.SetCellValue("Sheet1", fmt.Sprintf("I%d", i+2), *transaction.Recepient.BankName)
	}
	fileName := fmt.Sprintf("Account Statement %s.xlsx", utils.GenerateUUIDString())
	file.SaveAs(fileName)
	logger.Info("file generated and saved", logger.LoggerOptions{
		Key: "file name",
		Data: fileName,
	})
	return &fileName, nil
}