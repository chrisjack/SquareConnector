// services/OpenConnectService.js
export class OpenConnectService {
    constructor(apiClient, eventBus) {
        this.api = apiClient;
        this.eventBus = eventBus;
        this.status = null;
    }

    async connect(ip, port) {
        if (!ip || !port) {
            this.eventBus.emit('openconnect:error', 'Please enter valid IP and port');
            return { success: false };
        }

        try {
            const response = await this.api.post('/api/openconnect/connect', { ip, port });

            if (response.ok) {
                this.eventBus.emit('openconnect:connecting');
                return { success: true };
            } else {
                throw new Error(`Failed to connect: ${response.statusText}`);
            }
        } catch (error) {
            this.eventBus.emit('openconnect:error', error.message);
            return { success: false, error: error.message };
        }
    }

    async disconnect() {
        try {
            const response = await this.api.post('/api/openconnect/disconnect');

            if (response.ok) {
                this.eventBus.emit('openconnect:disconnecting');
                return { success: true };
            }
        } catch (error) {
            this.eventBus.emit('openconnect:error', error.message);
            return { success: false, error: error.message };
        }
    }

    async saveConfig(ip, port, autoConnect) {
        try {
            const response = await this.api.post('/api/openconnect/config', {
                ip, port, autoConnect
            });

            if (!response.ok) {
                throw new Error(`Failed to save config: ${response.statusText}`);
            }

            return { success: true };
        } catch (error) {
            this.eventBus.emit('openconnect:error', error.message);
            return { success: false, error: error.message };
        }
    }

    updateStatus(status) {
        this.status = status;
        this.eventBus.emit('openconnect:status', status);
    }
}
