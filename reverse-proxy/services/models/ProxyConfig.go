package models

import("time")

type ProxyConfig struct {
	Port int `json:"port"`
	Admin_port int `json:"admin_port"`
	Strategy string `json:"strategy"` // e.g., "round-robin" or "least-conn"
	HealthCheckFreq time.Duration `json:"-"`
	Timeout time.Duration `json:"-"`
}