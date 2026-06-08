package utils

import beego "github.com/beego/beego/v2/server/web"

// configOrDefault reads a string config key, returning fallback if missing or empty
func configOrDefault(key, fallback string) string {
    val, err := beego.AppConfig.String(key)
    if err != nil || val == "" {
        return fallback
    }
    return val
}