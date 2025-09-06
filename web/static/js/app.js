class AudioRecognitionApp {
    constructor() {
        this.socket = null;
        this.recording = false;
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.connectWebSocket();
    }

    setupEventListeners() {
        const recordBtn = document.getElementById('start-recording');
        if (recordBtn) {
            recordBtn.addEventListener('click', () => {
                this.startRecording();
            });
        }

        const addSongForm = document.getElementById('add-song-form');
        if (addSongForm) {
            addSongForm.addEventListener('submit', (e) => {
                e.preventDefault();
                this.addSong();
            });
        }
    }

    setButtonLoading(btn, isLoading, loadingText = 'Buffering...') {
        if (!btn) return;
        if (isLoading) {
            if (!btn.dataset.originalHtml) {
                btn.dataset.originalHtml = btn.innerHTML;
            }
            btn.disabled = true;
            btn.classList.add('is-loading');
            btn.innerHTML = `<i class="fas fa-spinner fa-spin"></i> ${loadingText}`;
        } else {
            if (btn.dataset.originalHtml) {
                btn.innerHTML = btn.dataset.originalHtml;
            }
            btn.disabled = false;
            btn.classList.remove('is-loading');
        }
    }

    connectWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws`;
        
        this.socket = new WebSocket(wsUrl);
        
        this.socket.onopen = () => {
            console.log('WebSocket connected');
        };
        
        this.socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            this.handleWebSocketMessage(data);
        };
        
        this.socket.onclose = () => {
            console.log('WebSocket disconnected');
            setTimeout(() => this.connectWebSocket(), 3000);
        };
    }

    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'recording_status':
                this.updateRecordingStatus(data.payload);
                break;
            case 'early_guess':
                const eg = document.getElementById('early-guess');
                if (eg && data.payload && data.payload.name) {
                    eg.textContent = data.payload.name;
                }
                break;
            case 'result':
                this.showResult(data.payload);
                break;
            case 'error':
                this.showError(data.payload.message);
                break;
        }
    }

    async startRecording() {
        if (this.recording) return;
        
        this.recording = true;
        this.showRecordingStatus();
        
        try {
            const response = await fetch('/api/record', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                }
            });
            
            if (!response.ok) {
                throw new Error('Failed to start recording');
            }
            
            const result = await response.json();
            console.log('Recording started:', result);
            
        } catch (error) {
            console.error('Recording error:', error);
            this.showError('Failed to start recording');
            this.recording = false;
        }
    }

    showRecordingStatus() {
        const statusEl = document.getElementById('recording-status');
        const actionsEl = document.querySelector('.main-actions');
        
        if (statusEl) statusEl.style.display = 'block';
        if (actionsEl) actionsEl.style.display = 'none';
        
        this.startCountdown();
    }

    startCountdown() {
        const countdownElement = document.getElementById('countdown');
        const countdownNumber = document.getElementById('countdown-number');
        
        if (!countdownElement || !countdownNumber) return;
        
        countdownElement.style.display = 'block';
        
        let count = 3;
        const interval = setInterval(() => {
            countdownNumber.textContent = count;
            count--;
            
            if (count < 0) {
                clearInterval(interval);
                countdownElement.style.display = 'none';
                this.updateRecordingStatus({
                    status: 'recording',
                    message: 'Recording now... Play your song!',
                    progress: 0
                });
            }
        }, 1000);
    }

    updateRecordingStatus(status) {
        const titleEl = document.getElementById('status-title');
        const messageEl = document.getElementById('status-message');
        const progressEl = document.getElementById('progress');
        
        if (titleEl) {
            titleEl.textContent = status.status === 'recording' ? 'Recording Audio...' : 'Processing...';
        }
        if (messageEl) {
            messageEl.textContent = status.message;
        }
        if (progressEl && status.progress !== undefined) {
            progressEl.style.width = `${status.progress}%`;
        }
    }

    showResult(result) {
        const statusEl = document.getElementById('recording-status');
        const resultsEl = document.getElementById('results');
        
        if (statusEl) statusEl.style.display = 'none';
        if (resultsEl) resultsEl.style.display = 'block';
        
        if (result.is_match && result.song) {
            this.displaySuccessResult(result);
        } else {
            //this.displaySuccessResult(result);
            this.displayNoMatchResult(result);
        }
        
        this.recording = false;
    }

    displaySuccessResult(result) {
        const elements = {
            icon: document.getElementById('result-icon'),
            title: document.getElementById('result-title'),
            songTitle: document.getElementById('song-title'),
            artist: document.getElementById('song-artist'),
            album: document.getElementById('song-album'),
            confidence: document.getElementById('confidence-percent'),
            fill: document.getElementById('confidence-fill')
        };
        
        if (elements.icon) elements.icon.className = 'fas fa-check-circle';
        if (elements.title) elements.title.textContent = 'Song Identified!';
        if (elements.songTitle) elements.songTitle.textContent = result.song.title;
        if (elements.artist) elements.artist.textContent = result.song.artist;
        if (elements.album) elements.album.textContent = result.song.album || 'Unknown Album';
        
        const confidence = Math.round(result.confidence * 100);
        if (elements.confidence) elements.confidence.textContent = confidence;
        if (elements.fill) {
            elements.fill.style.width = `${confidence}%`;
            elements.fill.style.backgroundColor = confidence >= 60 ? '#4CAF50' : 
                                                 confidence >= 40 ? '#FF9800' : '#F44336';
        }
    }

    displayNoMatchResult(result) {
        const elements = {
            icon: document.getElementById('result-icon'),
            title: document.getElementById('result-title'),
            songTitle: document.getElementById('song-title'),
            artist: document.getElementById('song-artist'),
            album: document.getElementById('song-album'),
            confidence: document.getElementById('confidence-percent'),
            fill: document.getElementById('confidence-fill')
        };
        
        if (elements.icon) elements.icon.className = 'fas fa-times-circle';
        if (elements.title) elements.title.textContent = 'Song Not Found';
        if (elements.songTitle) elements.songTitle.textContent = 'Unknown';
        if (elements.artist) elements.artist.textContent = 'This song is not in our database';
        if (elements.album) elements.album.textContent = '';
        
        const confidence = Math.round(result.confidence * 100);
        if (elements.confidence) elements.confidence.textContent = confidence;
        if (elements.fill) {
            elements.fill.style.width = `${confidence}%`;
            elements.fill.style.backgroundColor = '#F44336';
        }
    }

    showError(message) {
        const statusEl = document.getElementById('recording-status');
        const resultsEl = document.getElementById('results');
        
        if (statusEl) statusEl.style.display = 'none';
        if (resultsEl) resultsEl.style.display = 'block';
        
        const elements = {
            icon: document.getElementById('result-icon'),
            title: document.getElementById('result-title'),
            songTitle: document.getElementById('song-title'),
            artist: document.getElementById('song-artist')
        };
        
        if (elements.icon) elements.icon.className = 'fas fa-exclamation-triangle';
        if (elements.title) elements.title.textContent = 'Error';
        if (elements.songTitle) elements.songTitle.textContent = message;
        if (elements.artist) elements.artist.textContent = 'Please try again';
        
        this.recording = false;
    }

    async addSong() {
        const form = document.getElementById('add-song-form');
        if (!form) return;

        const submitBtn = form.querySelector('button[type="submit"]');
        const formControls = form.querySelectorAll('input, button');

        const formData = new FormData(form);
        const songData = {
            artist: formData.get('artist'),
            title: formData.get('title'),
            album: formData.get('album')
        };

        this.setButtonLoading(submitBtn, true, 'Buffering...');
        formControls.forEach(el => el.disabled = true);

        try {
            const response = await fetch('/api/songs/add', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(songData)
            });

            if (!response.ok) {
                throw new Error('Failed to add song');
            }

            const result = await response.json();
            console.log('Song added:', result);

            this.closeAddSongModal();
            this.showNotification('Song added successfully!', 'success');

            setTimeout(() => window.location.reload(), 1500);
        } catch (error) {
            console.error('Add song error:', error);
            this.showNotification('Failed to add song', 'error');
        } finally {
            this.setButtonLoading(submitBtn, false);
            formControls.forEach(el => el.disabled = false);
        }
    }

    resetApp() {
        const resultsEl = document.getElementById('results');
        const actionsEl = document.querySelector('.main-actions');
        const progressEl = document.getElementById('progress');
        
        if (resultsEl) resultsEl.style.display = 'none';
        if (actionsEl) actionsEl.style.display = 'block';
        if (progressEl) progressEl.style.width = '0%';
    }

    showNotification(message, type = 'info') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        setTimeout(() => notification.classList.add('show'), 100);
        
        setTimeout(() => {
            notification.classList.remove('show');
            setTimeout(() => {
                if (document.body.contains(notification)) {
                    document.body.removeChild(notification);
                }
            }, 300);
        }, 3000);
    }

    closeAddSongModal() {
        const modal = document.getElementById('add-song-modal');
        const form = document.getElementById('add-song-form');
        
        if (modal) modal.style.display = 'none';
        if (form) form.reset();
    }
}

function showAddSongModal() {
    const modal = document.getElementById('add-song-modal');
    if (modal) modal.style.display = 'block';
}

function closeAddSongModal() {
    if (window.app) {
        window.app.closeAddSongModal();
    }
}

function resetApp() {
    if (window.app) {
        window.app.resetApp();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    window.app = new AudioRecognitionApp();
});