package delivery

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"notification/domain"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
)

type studentParentHandler struct {
	uc domain.StudentParentUseCase
}

func NewStudentParentHandler(app *fiber.App, useCase domain.StudentParentUseCase) {
	handler := &studentParentHandler{
		uc: useCase,
	}

	route := app.Group("/student-and-parent")
	route.Post("/insert", handler.CreateStudentAndParent)
	route.Post("/import", handler.UploadAndImport)
	route.Put("/modify/:id", handler.UpdateStudentAndParent)
	route.Delete("/rm/:id", handler.DeleteStudentAndParent)
	route.Get("/student/:id", handler.GetStudentDetailsByID)
}

func (sph *studentParentHandler) CreateStudentAndParent(c *fiber.Ctx) error {
	var req domain.StudentAndParent
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	_, err := govalidator.ValidateStruct(req.Student)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid Student request body",
		})
	}

	_, err = govalidator.ValidateStruct(req.Parent)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid Parent request body",
		})
	}

	if err := sph.uc.CreateStudentAndParentUC(c.Context(), &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err,
			"message": "Failed to Create Student and Parent",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Student and Parent created successfully",
	})
}

func (sph *studentParentHandler) UploadAndImport(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to parse file",
		})
	}

	// Define upload directory
	uploadDir := "./uploads"
	// Ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	// Save the file
	filePath := filepath.Join(uploadDir, file.Filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to save file",
		})
	}

	// Process the CSV file and get duplicate records
	resDupe, invalidTelephones, err := sph.processCSVFile(c.Context(), filePath)

	if invalidTelephones != nil && err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":            false,
			"error":              err.Error(),
			"message":            "Failed to process CSV file",
			"invalid_telephones": invalidTelephones,
		})
	}

	if resDupe != nil && err == nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success":    false,
			"message":    "File processed successfully, but some duplicates were found.",
			"duplicates": resDupe,
		})
	}

	// If no errors and no duplicates, return success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "File processed successfully",
	})

}

func (sph *studentParentHandler) processCSVFile(c context.Context, filePath string) (*[]string, *[]string, error) {
	var listStudentAndParent []domain.StudentAndParent
	var parentDataHolder domain.Parent
	var studentDataHolder domain.Student

	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()
	defer func() {
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to delete file: %v", err)
		}
	}()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV file: %v", err)
	}

	// Start from row 2 karena row 1 header
	for i, row := range records[1:] {
		if len(row) < 8 {
			log.Printf("Skipping row %d due to insufficient columns", i+2)
			continue
		}

		studentDataHolder = domain.Student{
			Name:      row[0],
			Class:     row[1],
			Gender:    row[2],
			Telephone: row[3],
			ParentID:  0,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = govalidator.ValidateStruct(studentDataHolder)
		if err != nil {
			return nil, nil, fmt.Errorf("row %d: error validating student: %v", i+2, err)
		}

		parentDataHolder = domain.Parent{
			Name:      row[4],
			Gender:    row[5],
			Telephone: row[6],
			Email:     getStringPointer(row[7]),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = govalidator.ValidateStruct(parentDataHolder)
		if err != nil {
			log.Printf("Parent validation failed on row %d: %v", i+2, err)
			return nil, nil, fmt.Errorf("row %d: error validating parent: %v", i+2, err)
		}

		studNParent := domain.StudentAndParent{
			Student: studentDataHolder,
			Parent:  parentDataHolder,
		}
		// Append to the list
		listStudentAndParent = append(listStudentAndParent, studNParent)
	}

	duplicates, err := sph.uc.ImportCSV(c, &listStudentAndParent)
	if err != nil {
		return nil, nil, fmt.Errorf("error importing CSV data: %v", err)
	}

	if len(*duplicates) > 0 {
		return duplicates, nil, nil
	}

	return nil, nil, nil
}

// Helper function to get a pointer to a string
func getStringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (sph *studentParentHandler) UpdateStudentAndParent(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID Required",
		})
	}

	convertetID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid ID",
			"error":   err.Error(),
		})
	}

	// Check if student exists
	_, err = sph.uc.GetStudentDetailsByID(c.Context(), convertetID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Targeted Student and Parent doesn't exist",
		})
	}

	var req domain.StudentPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	fmt.Println(req)

	// Validate request body
	_, err = govalidator.ValidateStruct(&req)
	if err != nil {
		// Get validation errors as a map
		validationErrors := govalidator.ErrorsByField(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"errors":  validationErrors,
			"message": "Invalid request body",
		})
	}

	// Perform the update operation
	if err := sph.uc.UpdateStudentAndParent(c.Context(), convertetID, &req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
			"message": "Failed to update student and parent",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Student and Parent updated successfully",
	})
}

func (sph *studentParentHandler) DeleteStudentAndParent(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid student ID",
			"error":   err.Error(),
		})
	}

	if err := sph.uc.DeleteStudentAndParent(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to delete student",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Student deleted successfully",
	})
}

func (sph *studentParentHandler) GetStudentDetailsByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Invalid student ID",
			"error":   err.Error(),
		})
	}

	student, err := sph.uc.GetStudentDetailsByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to get student",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Student and Parent retrieved successfully",
		"data":    student.Student,
	})
}
