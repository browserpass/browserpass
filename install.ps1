# Installs the native manifest on windows
# 

$app = 'com.dannyvankooten.browserpass'

$dirpath = Join-Path -Path $env:localappdata -ChildPath 'browserpass'
$ff_jsonpath = Join-Path -Path $dirpath -ChildPath "$app-firefox.json"
$chrome_jsonpath = Join-Path -Path $dirpath -ChildPath "$app-chrome.json"

# Make our local directory
new-item -type Directory -Path $dirpath -force

# copy our bin to local directory
& Copy-Item browserpass-windows64.exe $dirpath

# copy the native messaging manifest
$ffile = gc firefox-host.json
$ffile -replace '%%replace%%', ((Join-Path -Path $dirpath -ChildPath 'browserpass-windows64.exe' | ConvertTo-json) -replace '^"|"$', "") | Out-File -Encoding UTF8 $ff_jsonpath

$cfile = gc chrome-host.json
$cfile -replace '%%replace%%', ((Join-Path -Path $dirpath -ChildPath 'browserpass-windows64.exe' | ConvertTo-json) -replace '^"|"$', "") | Out-File -Encoding UTF8 $chrome_jsonpath

Write-Host ""
Write-Host "Which browser are you using?"
Write-Host "1) Firefox"
Write-Host "2) Chrome"

$browser = Read-Host
$allUsers = Read-Host "Install for all users? (y/n) [n]"

$installDest = "cu" # Current User
switch -wildcard ($allUsers) {
    "y*" { $installDest = "lm"; Break; }
    "n*" { $installDest = "cu"; Break; }
    default { $installDest = "cu" }
}

$browserToUse = ""
switch -regex ($browser) {
    '^1$|^f' { $browserToUse = "Mozilla"; Break; }
    '^2$|^c' { $browserToUse = "Google/Chrome"; Break; }
    default {$browserToUse = "unknown"}
}

If ($installDest -eq "lm" -And -NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole(`
            [Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Warning "Please re-run this script with Admin rights!"
    exit
}

$regPath = "hk{0}:\Software\{1}\NativeMessagingHosts" -f $installDest, $browserToUse

If (-NOT (Test-Path -Path $regPath)) {
    New-Item -Path $regPath -force
}
New-Item -Path "$regPath\$app" -force
New-ItemProperty -Path "$regPath\$app"`
    -Name '(Default)'`
    -Value $(If ($browserToUse -eq "Mozilla") {$ff_jsonpath} Else {$chrome_jsonpath})