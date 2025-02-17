package domain

import (
	"time"
)

type SafeStaffData struct {
	UserID    int        `json:"user_id"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Teaching  []Subject  `json:"teaching"`
}

type SafeStaffUpdatePayload struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StuAndPar struct {
	Student StudentBase `json:"student"`
	Parent  ParentBase  `json:"parent"`
}

type StudentBase struct {
	StudentNSN string `gorm:"type:varchar(10);not null;" json:"student_nsn" valid:"required~NSN is required"`
	Name       string `gorm:"type:varchar(150);not null;" json:"name" valid:"required~Name is required"`
	Class      string `gorm:"type:varchar(3);not null" json:"class" valid:"required~Class is required"`
	Gender     string `gorm:"type:gender_enum;not null" json:"gender" valid:"required~Gender is required"`
	Telephone  string `gorm:"type:varchar(13);not null" json:"telephone" valid:"required~Telephone is required"`
	ParentID   int    `gorm:"not null" json:"parent_id"`
}

type ParentBase struct {
	ParentID  int        `json:"parent_id" valid:"required~Parent ID is required"`
	Name      string     `json:"name" valid:"required~Name is required"`
	Gender    string     `json:"gender" valid:"required~Gender is required,in(male|female)~Invalid gender"`
	Telephone string     `json:"telephone" valid:"required~Telephone is required"`
	Email     *string    `json:"email,omitempty" valid:"email~Invalid email format"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type StudentAndParentBase struct {
	StudentNSN string     `json:"student_nsn" valid:"required~NSN is required"`
	Name       string     `json:"name" valid:"required~Name is required"`
	Class      string     `json:"class" valid:"required~Class is required"`
	Gender     string     `json:"gender" valid:"required~Gender is required,in(male|female)~Invalid gender"`
	Telephone  string     `json:"telephone" valid:"required~Telephone is required"`
	ParentID   int        `json:"parent_id"`
	Parent     ParentBase `json:"parent" valid:"required~Parent details are required"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

type StudentPayload struct {
	Student Student `json:"student"`
}

type StudentAndParent2 struct {
	Student StudentNoGorm `json:"student"`
	Parent  ParentNoGorm  `json:"parent"`
}

type ParentNoGorm struct {
	ParentID  int        `json:"parent_id"`
	Name      string     `json:"name" valid:"required~Name is required"`
	Gender    string     `json:"gender" valid:"required~Gender is required,in(male|female)~Invalid gender"`
	Telephone string     `json:"telephone" valid:"required~Telephone is required"`
	Email     *string    `json:"email" valid:"email~Invalid email format,optional~true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type StudentNoGorm struct {
	StudentID int        `json:"student_id"`
	NSN       string     `gorm:"type:varchar(10);not null;" json:"student_nsn" valid:"required~NSN is required"`
	Name      string     `json:"name" valid:"required~Name is required"`
	Class     string     `json:"class" valid:"required~Class is required"`
	Gender    string     `json:"gender" valid:"required~Gender is required,in(male|female)~Invalid gender"`
	Telephone string     `json:"telephone" valid:"required~Telephone is required"`
	ParentID  int        `json:"parent_id"`
	Parent    Parent     `json:"parent"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type UserResponse struct {
	UserID    int        `json:"user_id"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type AttendanceNotificationHistoryResponse struct {
	Student        Student      `json:"student"`
	Parent         Parent       `json:"parent"`
	User           UserResponse `json:"user"`
	Subject        Subject      `json:"subject"`
	WhatsappStatus bool         `json:"whatsapp_status"`
	EmailStatus    bool         `json:"email_status"`
	CreatedAt      time.Time    `json:"created_at"`
}

type StudentTestScore struct {
	StudentNSN string   `json:"student_nsn"`
	TestScore  *float64 `json:"test_score"`
}

type InputTestScorePayload struct {
	StudentTestScore []StudentTestScore `json:"students_test_score"`
	SubjectCode      string             `json:"subject_code"`
}

type SubjectAndScoreResult struct {
	SubjectCode string   `json:"subject_code"`
	Subject     Subject  `json:"subject"`
	Score       *float64 `json:"score"`
}

type IndividualExamScore struct {
	StudentNSN            string                  `json:"student_nsn"`
	Student               Student                 `json:"student"`
	SubjectAndScoreResult []SubjectAndScoreResult `json:"subject_and_score_result"`
}

type StudentsAssociateWithParent struct {
	Parent            Parent    `json:"parent"`
	AssociatedStudent []Student `json:"students"`
}
