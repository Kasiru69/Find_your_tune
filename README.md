# Find_your_tune
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Find Your Tune - Installation Guide</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .step { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .code { background: #2d2d2d; color: #f8f8f8; padding: 10px; border-radius: 5px; font-family: monospace; }
        .warning { background: #fff3cd; border: 1px solid #ffeaa7; padding: 10px; border-radius: 5px; }
        h2 { color: #2c3e50; }
        h3 { color: #e74c3c; }
    </style>
</head>
<body>

<h1>🎵 Find Your Tune - Installation Guide</h1>

<h2>📋 Prerequisites</h2>
<div class="step">
    <h3>Required Software:</h3>
    <ul>
        <li><strong>Go 1.21+</strong> - <a href="https://golang.org/dl/">Download here</a></li>
        <li><strong>Git</strong> - <a href="https://git-scm.com/downloads">Download here</a></li>
        <li><strong>GCC Compiler</strong> (MinGW for Windows)</li>
        <li><strong>FFmpeg</strong> - <a href="https://ffmpeg.org/download.html">Download here</a></li>
        <li><strong>yt-dlp</strong> - <code>pip install yt-dlp</code></li>
    </ul>
</div>

<h2>🔧 Windows Setup</h2>

<div class="step">
    <h3>Step 1: Install MinGW-w64 (GCC for Windows)</h3>
    <ol>
        <li>Download MinGW from <a href="https://winlibs.com/">winlibs.com</a></li>
        <li>Extract to <code>C:\mingw64</code></li>
        <li>Add <code>C:\mingw64\bin</code> to System PATH</li>
        <li>Restart your terminal/IDE</li>
    </ol>
    
    <h4>Verify Installation:</h4>
    <div class="code">gcc --version</div>
</div>

<div class="step">
    <h3>Step 2: Enable CGO Permanently</h3>
    <div class="code">go env -w CGO_ENABLED=1</div>
    
    <h4>Verify CGO is enabled:</h4>
    <div class="code">go env CGO_ENABLED</div>
    <p><em>Should show: 1</em></p>
</div>

<h2>🚀 Installation & Running</h2>

<div class="step">
    <h3>Step 1: Clone the Repository</h3>
    <div class="code">
git clone https://github.com/Kasiru69/Find_your_tune.git<br>
cd Find_your_tune
    </div>
</div>

<div class="step">
    <h3>Step 2: Install Go Dependencies</h3>
    <div class="code">go mod tidy</div>
</div>

<div class="step">
    <h3>Step 3: Run the Application</h3>
    <div class="code">go run ./cmd/server</div>
    
    <div class="warning">
        <strong>Windows Users:</strong> If you get CGO errors, run:
        <div class="code">
$env:CGO_ENABLED=1<br>
go clean -cache<br>
go run ./cmd/server
        </div>
    </div>
</div>

<div class="step">
    <h3>Step 4: Access the Web Interface</h3>
    <p>Open your browser and go to:</p>
    <div class="code">http://localhost:8080</div>
</div>

<h2>🐧 Linux/macOS Setup</h2>

<div class="step">
    <h3>Install Dependencies</h3>
    
    <h4>Ubuntu/Debian:</h4>
    <div class="code">
sudo apt-get update<br>
sudo apt-get install build-essential ffmpeg<br>
pip install yt-dlp
    </div>
    
    <h4>macOS:</h4>
    <div class="code">
brew install gcc ffmpeg<br>
pip install yt-dlp
    </div>
    
    <h4>Enable CGO and Run:</h4>
    <div class="code">
export CGO_ENABLED=1<br>
go mod tidy<br>
go run ./cmd/server
    </div>
</div>

<h2>✅ Expected Output</h2>
<div class="step">
    <p>When successfully running, you should see:</p>
    <div class="code">
🎵 Audio Recognition Server starting on http://localhost:8080<br>
📁 Database: data/songs.db<br>
📁 Temp dir: data/temp
    </div>
</div>

<h2>🎯 Usage</h2>
<div class="step">
    <ol>
        <li><strong>Add Songs:</strong> Use "Add New Song" to build your database</li>
        <li><strong>Record Audio:</strong> Click "Start Recording" and play any song</li>
        <li><strong>View Results:</strong> Get match results with confidence percentage</li>
        <li><strong>Manage Database:</strong> Browse and search your song collection</li>
    </ol>
</div>

<h2>🔧 Troubleshooting</h2>

<div class="step">
    <h3>Common Issues:</h3>
    
    <h4>CGO Error:</h4>
    <div class="code">
# Windows
$env:CGO_ENABLED=1

# Linux/macOS  
export CGO_ENABLED=1
    </div>
    
    <h4>GCC Not Found:</h4>
    <p>Install MinGW-w64 and add to PATH: <code>C:\mingw64\bin</code></p>
    
    <h4>FFmpeg Not Found:</h4>
    <p>Download FFmpeg and add to system PATH</p>
    
    <h4>Port Already in Use:</h4>
    <p>Change port in <code>config/config.go</code> or kill existing process</p>
</div>

<div class="warning">
    <h3>⚠️ Important Notes</h3>
    <ul>
        <li>Ensure all dependencies are in your system PATH</li>
        <li>Restart terminal after PATH changes</li>
        <li>Windows users need MinGW-w64 for CGO support</li>
        <li>The application creates database and temp directories automatically</li>
    </ul>
</div>

<h2>📞 Support</h2>
<div class="step">
    <p>If you encounter issues:</p>
    <ul>
        <li>Check all prerequisites are installed</li>
        <li>Verify PATH environment variables</li>
        <li>Run commands in a fresh terminal</li>
        <li>Create an issue on GitHub for help</li>
    </ul>
</div>

</body>
</html>
