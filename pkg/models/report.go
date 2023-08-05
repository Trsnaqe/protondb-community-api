package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`
	Data          map[string]interface{} `bson:"data"`
	ReportVersion string                 `bson:"report_version"`
}

type ReportFormatV2 struct {
	App struct {
		Steam struct {
			AppID interface{} `json:"appId,omitempty"`
		} `json:"steam"`
		Title string `json:"title,omitempty"`
	} `json:"app,omitempty"`
	Responses  map[string]interface{} `json:"responses,omitempty"`
	Timestamp  int64                  `json:"timestamp,omitempty"`
	SystemInfo struct {
		CPU       string `json:"cpu,omitempty"`
		GPU       string `json:"gpu,omitempty"`
		GPUDriver string `json:"gpuDriver,omitempty"`
		Kernel    string `json:"kernel,omitempty"`
		OS        string `json:"os,omitempty"`
		RAM       string `json:"ram,omitempty"`
	} `json:"systemInfo,omitempty"`
}

type ReportFormatV1 struct {
	AppID         string                  `json:"appId,omitempty"`
	Title         string                  `json:"title,omitempty"`
	Timestamp     interface{}             `json:"timestamp,omitempty"`
	Rating        interface{}             `json:"rating,omitempty"`
	OS            interface{}             `json:"os,omitempty"`
	Notes         *interface{}            `json:"notes,omitempty"`
	GPUDriver     *interface{}            `json:"gpuDriver,omitempty"`
	Specs         *interface{}            `json:"specs,omitempty"`
	ProtonVersion *interface{}            `json:"protonVersion,omitempty"`
	CPU           *interface{}            `json:"cpu,omitempty"`
	Duration      *interface{}            `json:"duration,omitempty"`
	GPU           *interface{}            `json:"gpu,omitempty"`
	Kernel        *interface{}            `json:"kernel,omitempty"`
	RAM           *interface{}            `json:"ram,omitempty"`
	Tweaks        *map[string]interface{} `json:"tweaks,omitempty"`
}
