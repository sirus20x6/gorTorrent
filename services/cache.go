cache, err := services.NewCache("plugins")
if err != nil {
    log.Fatalf("Failed to create cache: %v", err)
}

// Store settings
settings := &CachedSettings{
    Modified: true,
    Hash: "abc123",
}
if err := cache.Set("settings.dat", settings); err != nil {
    log.Printf("Failed to cache settings: %v", err)
}

// Retrieve settings
var loaded CachedSettings
if err := cache.Get("settings.dat", &loaded); err != nil {
    log.Printf("Failed to load settings: %v", err)
}