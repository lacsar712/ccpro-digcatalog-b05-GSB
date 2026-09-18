package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"uniqueIndex;size:64;not null"`
	PasswordHash string         `json:"-" gorm:"size:255;not null"`
	Role         string         `json:"role" gorm:"size:32;not null"` // admin | recorder
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

type Site struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:128;not null"`
	Period    string         `json:"period" gorm:"size:64;not null"` // 新石器/商周等
	Latitude  float64        `json:"latitude"`
	Longitude float64        `json:"longitude"`
	Manager   string         `json:"manager" gorm:"size:64"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Units     []Unit         `json:"units,omitempty" gorm:"foreignKey:SiteID"`
}

type Unit struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	SiteID           uint           `json:"siteId" gorm:"not null;index"`
	Code             string         `json:"code" gorm:"size:64;not null"` // T1, T2...
	DepthMin         float64        `json:"depthMin"`
	DepthMax         float64        `json:"depthMax"`
	StratumDesc      string         `json:"stratumDesc" gorm:"type:text"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
	Site             *Site          `json:"site,omitempty" gorm:"foreignKey:SiteID"`
	Finds            []Find         `json:"finds,omitempty" gorm:"foreignKey:UnitID"`
}

type Material struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:64;not null"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// FindCompletenessLog 记录 Find.completeness 的每次变更轨迹（仅在实际变化时写入）。
type FindCompletenessLog struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	FindID     uint      `json:"findId" gorm:"not null;index"`
	FromValue  string    `json:"fromValue" gorm:"size:32"`
	ToValue    string    `json:"toValue" gorm:"size:32"`
	OperatorID uint      `json:"operatorId" gorm:"not null;index"`
	ChangedAt  time.Time `json:"changedAt" gorm:"not null;index"`
	Note       *string   `json:"note" gorm:"type:text"` // 可空备注
	Operator   *User     `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

type Find struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UnitID       uint           `json:"unitId" gorm:"not null;index"`
	MaterialID   *uint          `json:"materialId" gorm:"index"`
	RegisterNo   string         `json:"registerNo" gorm:"uniqueIndex;size:64;not null"`
	ArtifactType string         `json:"artifactType" gorm:"size:64;not null"` // 陶片/青铜器/骨器
	MaterialName string         `json:"materialName" gorm:"size:64"`          // 冗余展示字段
	Completeness string         `json:"completeness" gorm:"size:32"`          // 完整/残缺/碎片
	FindDate     *time.Time     `json:"findDate" gorm:"type:date"`
	Description  string         `json:"description" gorm:"type:text"`
	StorageLoc   string         `json:"storageLoc" gorm:"size:128"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Unit         *Unit          `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
	Material     *Material      `json:"material,omitempty" gorm:"foreignKey:MaterialID"`
}
