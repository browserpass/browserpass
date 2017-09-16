# Installs the native manifest on windows
# 

$app = 'com.dannyvankooten.browserpass'

$dirpath = Join-Path -Path $env:localappdata -ChildPath 'browserpass'
$ff_jsonpath = Join-Path -Path $dirpath -ChildPath "$app-firefox.json"
$chrome_jsonpath = Join-Path -Path $dirpath -ChildPath "$app-chrome.json"

# Make our local directory
new-item -type Directory -Path $dirpath -force

# copy our bin to local directory
& cp browserpass-windows64.exe $dirpath

# copy the native messaging manifest
$ffile = gc firefox-host.json
$ffile -replace '%%replace%%', ((Join-Path -Path $dirpath -ChildPath 'browserpass-windows64.exe' | ConvertTo-json) -replace '^"|"$', "") | Out-File -Encoding UTF8 $ff_jsonpath

$cfile = gc chrome-host.json
$cfile -replace '%%replace%%', ((Join-Path -Path $dirpath -ChildPath 'browserpass-windows64.exe' | ConvertTo-json) -replace '^"|"$', "") | Out-File -Encoding UTF8 $chrome_jsonpath

# Add oour registry values fore ff
New-Item -Path "hkcu:\Software\Mozilla\NativeMessagingHosts\$app" -force
New-ItemProperty -Path "hkcu:\Software\Mozilla\NativeMessagingHosts\$app" -Name '(Default)' -Value $ff_jsonpath -force

# Add oour registry values fore chrome
New-Item -Path "hkcu:\Software\Google\Chrome\NativeMessagingHosts\$app" -force
New-ItemProperty -Path "hkcu:\Software\Google\Chrome\NativeMessagingHosts\$app" -Name '(Default)' -Value $chrome_jsonpath -force

