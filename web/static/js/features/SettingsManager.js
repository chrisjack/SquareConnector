// features/SettingsManager.js
export class SettingsManager {
    constructor(apiClient, eventBus) {
        this.api = apiClient;
        this.eventBus = eventBus;
        this.settings = {};
    }

    async load() {
        try {
            const response = await this.api.get('/api/settings');
            if (response.ok) {
                this.settings = await response.json();
                this.eventBus.emit('settings:loaded', this.settings);
            }
        } catch (error) {
            console.error('Failed to load settings:', error);
        }
    }

    async save(newSettings) {
        try {
            // Merge with the last loaded settings before POSTing. The
            // server-side handleSettings handler unconditionally calls
            // every per-field setter, so fields omitted from the
            // payload would be persisted as their JSON zero values
            // (e.g. GSProIP would be wiped to ""). Sending a full
            // object preserves all unrelated fields.
            const merged = { ...this.settings, ...newSettings };
            const response = await this.api.post('/api/settings', merged);

            if (response.ok) {
                this.settings = merged;
                this.eventBus.emit('settings:saved', this.settings);
                return { success: true };
            } else {
                throw new Error(`Failed to save settings: ${response.statusText}`);
            }
        } catch (error) {
            this.eventBus.emit('settings:error', error.message);
            return { success: false, error: error.message };
        }
    }

    get(key) {
        return this.settings[key];
    }

    getAll() {
        return { ...this.settings };
    }
}
